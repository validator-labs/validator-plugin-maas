/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package controller defines a controller for reconciling MaasValidator objects.
package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	maasclient "github.com/canonical/gomaasclient/client"
	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	dnsval "github.com/validator-labs/validator-plugin-maas/internal/validators/dns"
	osval "github.com/validator-labs/validator-plugin-maas/internal/validators/os"
	resval "github.com/validator-labs/validator-plugin-maas/internal/validators/res"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
	vres "github.com/validator-labs/validator/pkg/validationresult"
)

// MaasValidatorReconciler reconciles a MaasValidator object
type MaasValidatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=maasvalidators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=maasvalidators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=maasvalidators/finalizers,verbs=update

// SetUpClient is defined to enable monkey patching the setUpClient function in integration tests
var SetUpClient = setUpClient

// Reconcile reconciles each rule found in each MaasValidator in the cluster and creates ValidationResults accordingly
func (r *MaasValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.V(0).WithValues("name", req.Name, "namespace", req.Namespace)
	l.Info("Reconciling MaasValidator")

	validator := &v1alpha1.MaasValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	secretName := validator.Spec.Auth.SecretName
	tokenKey := validator.Spec.Auth.TokenKey

	var (
		maasToken string
		err       error
	)

	if maasToken, err = r.tokenFromSecret(secretName, req.Namespace, tokenKey); err != nil {
		l.Error(err, "failed to retrieve MAAS API token")
	}

	maasURL := validator.Spec.Host
	maasClient, err := SetUpClient(maasURL, maasToken)

	if err != nil {
		l.Error(err, "failed to initialize MAAS client")
	}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	p, err := patch.NewHelper(vr, r.Client)
	if err != nil {
		l.Error(err, "failed to create patch helper")
		return ctrl.Result{}, err
	}
	nn := ktypes.NamespacedName{
		Name:      validationResultName(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExistingValidationResult(vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			l.Error(err, "unexpected error getting ValidationResult")
		}
		if err := vres.HandleNewValidationResult(ctx, r.Client, p, buildValidationResult(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Millisecond}, nil
	}

	// Always update the expected result count in case the validator's rules have changed
	vr.Spec.ExpectedResults = validator.Spec.ResultCount()

	resp := types.ValidationResponse{
		ValidationRuleResults: make([]*types.ValidationRuleResult, 0, vr.Spec.ExpectedResults),
		ValidationRuleErrors:  make([]error, 0, vr.Spec.ExpectedResults),
	}

	imageRulesService := osval.NewImageRulesService(r.Log, maasClient.BootResources)
	resourceRulesService := resval.NewResourceRulesService(r.Log, maasClient.Machines)
	upstreamDNSRulesService := dnsval.NewUpstreamDNSRulesService(r.Log, maasClient.MAASServer)
	internalDNSRulesService := dnsval.NewInternalDNSRulesService(r.Log, maasClient.DNSResources)

	// MAAS Instance image rules
	for _, rule := range validator.Spec.ImageRules {
		vrr, err := imageRulesService.ReconcileMaasInstanceImageRule(rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile MAAS image rule")
		}
		resp.AddResult(vrr, err)
	}

	// MAAS Instance upstream DNS rules
	for _, rule := range validator.Spec.UpstreamDNSRules {
		vrr, err := upstreamDNSRulesService.ReconcileMaasInstanceUpstreamDNSRule(rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile MAAS upstream DNS rule")
		}
		resp.AddResult(vrr, err)
	}

	seenAZ := make(map[string]bool, 0)
	// MAAS Instance resource availability rules
	for _, rule := range validator.Spec.ResourceAvailabilityRules {
		vrr, err := resourceRulesService.ReconcileMaasInstanceResourceRule(rule, seenAZ)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile MAAS resource rule")
		}
		resp.AddResult(vrr, err)
		seenAZ[rule.AZ] = true
	}

	// MAAS Instance internal DNS rules
	for _, rule := range validator.Spec.InternalDNSRules {
		vrr, err := internalDNSRulesService.ReconcileMaasInstanceInternalDNSRule(rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile MAAS internal DNS rule")
		}
		resp.AddResult(vrr, err)
	}

	// Patch the ValidationResult with the latest ValidationRuleResults
	if err := vres.SafeUpdateValidationResult(ctx, p, vr, resp, r.Log); err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Requeuing for re-validation in two minutes.")
	return ctrl.Result{RequeueAfter: time.Second * 120}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MaasValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MaasValidator{}).
		Complete(r)
}

func buildValidationResult(validator *v1alpha1.MaasValidator) *vapi.ValidationResult {
	return &vapi.ValidationResult{
		ObjectMeta: metav1.ObjectMeta{
			Name:      validationResultName(validator),
			Namespace: validator.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: validator.APIVersion,
					Kind:       validator.Kind,
					Name:       validator.Name,
					UID:        validator.UID,
					Controller: util.Ptr(true),
				},
			},
		},
		Spec: vapi.ValidationResultSpec{
			Plugin:          constants.PluginCode,
			ExpectedResults: validator.Spec.ResultCount(),
		},
	}
}

func validationResultName(validator *v1alpha1.MaasValidator) string {
	return fmt.Sprintf("validator-plugin-maas-%s", validator.Name)
}

func setUpClient(maasURL, maasToken string) (*maasclient.Client, error) {
	maasClient, err := maasclient.GetClient(maasURL, maasToken, "2.0")
	if err != nil {
		return nil, err
	}
	return maasClient, nil
}

func (r *MaasValidatorReconciler) tokenFromSecret(name, namespace, tokenKey string) (string, error) {
	r.Log.Info("Getting MAAS API token from secret", "name", name, "namespace", namespace)

	nn := ktypes.NamespacedName{Name: name, Namespace: namespace}
	secret := &corev1.Secret{}
	if err := r.Get(context.Background(), nn, secret); err != nil {
		return "", err
	}

	if key, found := secret.Data[tokenKey]; found {
		token := string(key)
		token = strings.TrimSuffix(token, "\n")
		return token, nil
	}
	return "", fmt.Errorf("secret does not contain %s", tokenKey)
}
