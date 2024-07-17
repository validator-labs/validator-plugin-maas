package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	// . "github.com/onsi/gomega"
	//+kubebuilder:scaffold:imports
	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
)

const MaasValidatorName = "maas-validator"

var _ = Describe("MaaSValidator controller", Ordered, func() {

	BeforeEach(func() {
		// toggle true/false to enable/disable the OCIValidator controller specs
		if false {
			Skip("skipping")
		}
	})

	val := &v1alpha1.MaasValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MaasValidatorName,
			Namespace: validatorNamespace,
		},
		Spec: v1alpha1.MaasValidatorSpec{
			Host: "maas.sc",
			Auth: v1alpha1.Auth{
				SecretName: "maas-api-token",
			},
			OSImageRules: []v1alpha1.OSImageRule{
				{Name: "Ubuntu", OSImages: []v1alpha1.OSImage{
					{OSName: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
			},
			InternalDNSRules: []v1alpha1.InternalDNSRule{
				{MaasDomain: "maas.sc", DNSRecords: []v1alpha1.DNSRecord{
					{Hostname: "maas.sc", Type: "A", IP: "10.0.0.1", TTL: 3600},
				}},
			},
			ExternalDNSRules: []v1alpha1.ExternalDNSRule{
				{Name: "Upstream DNS", Enabled: true},
			},
			ResourceAvailabilityRules: []v1alpha1.ResourceAvailabilityRule{
				{Name: "az1 2 machines", Resources: []v1alpha1.Resource{
					{AZ: "az1", NumMachines: 2, NumCPU: 2, NumDisk: 20, NumRAM: 4},
				}},
			},
		},
	}

	//secret := &corev1.Secret{}

	vr := &vapi.ValidationResult{}
	vrKey := types.NamespacedName{Name: validationResultName(val), Namespace: validatorNamespace}

	It("Should create a ValidationResult and update its Status with a failed condition", func() {
		By("By creating a new MaasValidator")
		ctx := context.Background()

		Expect(k8sClient.Create(ctx, val)).Should(Succeed())

		// Wait for the ValidationResult's Status to be updated
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, vrKey, vr); err != nil {
				return false
			}

			stateOk := vr.Status.State == vapi.ValidationFailed
			return stateOk
		}, timeout, interval).Should(BeTrue(), "failed to create a ValidationResult")
	})
})
