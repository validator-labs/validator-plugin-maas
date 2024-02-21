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

package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	maasclient "github.com/maas/gomaasclient/client"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ktypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-maas/internal/constants"
	val "github.com/spectrocloud-labs/validator-plugin-maas/internal/validators"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/util"
	vres "github.com/spectrocloud-labs/validator/pkg/validationresult"
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

// Reconcile reconciles each rule found in each OCIValidator in the cluster and creates ValidationResults accordingly
func (r *MaasValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.V(0).Info("Reconciling MaasValidator", "name", req.Name, "namespace", req.Namespace)

	validator := &v1alpha1.MaasValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	secretName := validator.Spec.MaasInstance.Auth.SecretName
	var (
		maasToken string = ""
		err       error  = nil
	)

	if maasToken, err = r.tokenFromSecret(secretName, req.Namespace); err != nil {
		r.Log.Error(err, "failed to retrieve MaaS API Key", "key", req)
	}

	maasUrl := validator.Spec.MaasInstance.Host
	maasclient, err := maasclient.GetClient(maasUrl, maasToken, "2.0")

	if err != nil {
		r.Log.Error(err, "failed to initialize MaaS client")
	}

	apiclient := val.MaaSAPI{Client: maasclient}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	nn := ktypes.NamespacedName{
		Name:      validationResultName(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExistingValidationResult(nn, vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			r.Log.V(0).Error(err, "unexpected error getting ValidationResult", "name", nn.Name, "namespace", nn.Namespace)
		}
		if err := vres.HandleNewValidationResult(r.Client, buildValidationResult(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
	}
	// Maas Instance image rules
	maasRuleService := val.NewMaasRuleService(&apiclient)
	validationResult, err := maasRuleService.ReconcileMaasInstanceImageRules(validator.Spec.MaasInstanceRules)
	if err != nil {
		r.Log.V(0).Error(err, "failed to reconcile MaaS instance rule")
	}
	vres.SafeUpdateValidationResult(r.Client, nn, validationResult, validator.Spec.ResultCount(), err, r.Log)

	// Maas Instance Ext DNS rules
	validationResult, err = maasRuleService.ReconsileMaasInstanceExtDNSRules(validator.Spec.MaasExternalDNSRule)
	if err != nil {
		r.Log.V(0).Error(err, "failed to reconcile MaaS instance rule for external DNS")
	}

	vres.SafeUpdateValidationResult(r.Client, nn, validationResult, validator.Spec.ResultCount(), err, r.Log)

	r.Log.V(0).Info("Requeuing for re-validation in two minutes.", "name", req.Name, "namespace", req.Namespace)
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

func (r *MaasValidatorReconciler) tokenFromSecret(name, namespace string) (string, error) {
	r.Log.Info("Getting MaaS API token from secret", "name", name, "namespace", namespace)

	nn := ktypes.NamespacedName{Name: name, Namespace: namespace}
	secret := &corev1.Secret{}
	if err := r.Get(context.Background(), nn, secret); err != nil {
		return "", err
	}

	if key, found := secret.Data["MAAS_API_KEY"]; found {
		token := string(key)
		token = strings.TrimSuffix(token, "\n")
		decodedBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return "", fmt.Errorf("failed to decode secret data")
		}
		return string(decodedBytes), nil
	}
	return "", fmt.Errorf("secret does not contain MAAS_API_KEY")
}
