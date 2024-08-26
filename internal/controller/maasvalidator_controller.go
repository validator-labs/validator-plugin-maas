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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/pkg/validate"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vres "github.com/validator-labs/validator/pkg/validationresult"
)

var errCredsRequired = errors.New("auth.secretName or auth.apiToken must be set")

// MaasValidatorReconciler reconciles a MaasValidator object
type MaasValidatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=maasvalidators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=maasvalidators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=maasvalidators/finalizers,verbs=update

// Reconcile reconciles each rule found in each MaasValidator in the cluster and creates ValidationResults accordingly
func (r *MaasValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.V(0).WithValues("name", req.Name, "namespace", req.Namespace)
	l.Info("Reconciling MaasValidator")

	validator := &v1alpha1.MaasValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	p, err := patch.NewHelper(vr, r.Client)
	if err != nil {
		l.Error(err, "failed to create patch helper")
		return ctrl.Result{}, err
	}
	nn := ktypes.NamespacedName{
		Name:      vres.Name(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExisting(vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			l.Error(err, "unexpected error getting ValidationResult")
		}
		if err := vres.HandleNew(ctx, r.Client, p, vres.Build(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Millisecond}, nil
	}

	// Always update the expected result count in case the validator's rules have changed
	vr.Spec.ExpectedResults = validator.Spec.ResultCount()

	if validator.Spec.Auth.APIToken == "" && validator.Spec.Auth.SecretName == "" {
		l.Error(errCredsRequired, "failed to reconcile MAAS validator with empty credentials")
		return ctrl.Result{}, errCredsRequired
	}

	if validator.Spec.Auth.SecretName != "" {
		validator.Spec.Auth.APIToken, err = r.tokenFromSecret(validator.Spec.Auth.SecretName, req.Namespace, validator.Spec.Auth.TokenKey)
		if err != nil {
			l.Error(err, "failed to retrieve MAAS API token from secret")
			return ctrl.Result{}, err
		}
	}

	resp := validate.Validate(validator.Spec, r.Log)

	// Patch the ValidationResult with the latest ValidationRuleResults
	if err := vres.SafeUpdate(ctx, p, vr, resp, r.Log); err != nil {
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
