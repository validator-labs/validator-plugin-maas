package controller

import (
	"context"

	"github.com/canonical/gomaasclient/api"
	maasclient "github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	// . "github.com/onsi/gomega"
	// +kubebuilder:scaffold:imports
	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/pkg/validate"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vres "github.com/validator-labs/validator/pkg/validationresult"
)

const (
	MaasValidatorName = "maas-validator"
	fakeToken         = "fake:api:token"
	fakeTokenKey      = "MAAS_API_KEY"
	secretName        = "maas-creds"
)

type MockBootResourcesService struct {
	api.BootResources
}

type MockUDNSRulesService struct {
	api.MAASServer
}

type MockMachinesService struct {
	api.Machines
}

type MockIDNSRulesService struct {
	api.DNSResources
}

func (b *MockBootResourcesService) Get(params *entity.BootResourcesReadParams) ([]entity.BootResource, error) {
	return []entity.BootResource{
		{
			Type:         "Synced",
			Name:         "Ubuntu",
			Architecture: "amd64/ga-22.04",
		},
	}, nil
}

func (u *MockUDNSRulesService) Get(string) ([]byte, error) {
	return []byte("8.8.8.8"), nil
}

func (m *MockMachinesService) Get(params *entity.MachinesParams) ([]entity.Machine, error) {
	return []entity.Machine{
		{
			Hostname: "maas.sc",
			CPUCount: 24,
			Memory:   128 * 1024,
			Storage:  1024 * 1000,
			Zone:     entity.Zone{Name: "az1"},
			Pool:     entity.ResourcePool{Name: "pool1"},
			TagNames: []string{"foo", "bar"},
		},
	}, nil
}

func (i *MockIDNSRulesService) Get(params *entity.DNSResourcesParams) ([]entity.DNSResource, error) {
	return []entity.DNSResource{
		{
			FQDN: "foo.maas.sc",
			ResourceRecords: []entity.DNSResourceRecord{
				{
					RRType: "A",
					RRData: "8.8.8.8",
					TTL:    10,
					ID:     0,
				},
			},
		},
	}, nil
}

var _ = Describe("MaaSValidator controller", Ordered, func() {

	BeforeEach(func() {
		// toggle true/false to enable/disable the OCIValidator controller specs
		if false {
			Skip("skipping")
		}
		// overwrite the maas client to inject mock services
		validate.SetUpClient = func(maasURL, massToken string) (*maasclient.Client, error) {
			c := &maasclient.Client{}
			c.BootResources = &MockBootResourcesService{}
			c.MAASServer = &MockUDNSRulesService{}
			c.Machines = &MockMachinesService{}
			c.DNSResources = &MockIDNSRulesService{}
			return c, nil
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
				SecretName: "maas-creds",
				TokenKey:   "MAAS_API_KEY",
			},
			ImageRules: []v1alpha1.ImageRule{
				{RuleName: "Ubuntu", Images: []v1alpha1.Image{
					{Name: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
			},
			InternalDNSRules: []v1alpha1.InternalDNSRule{
				{MaasDomain: "maas.sc", DNSResources: []v1alpha1.DNSResource{
					{FQDN: "maas.sc", DNSRecords: []v1alpha1.DNSRecord{
						{Type: "A", IP: "10.0.0.1", TTL: 3600},
					}},
				}},
			},
			UpstreamDNSRules: []v1alpha1.UpstreamDNSRule{
				{RuleName: "Upstream DNS", NumDNSServers: 1},
			},
			ResourceAvailabilityRules: []v1alpha1.ResourceAvailabilityRule{
				{RuleName: "az1 2 machines", AZ: "az1", Resources: []v1alpha1.Resource{
					{NumMachines: 2, NumCPU: 2, Disk: 20, RAM: 4},
				}},
			},
		},
	}

	vr := &vapi.ValidationResult{}
	vrKey := types.NamespacedName{Name: vres.Name(val), Namespace: validatorNamespace}

	maasSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: validatorNamespace,
		},
		Data: map[string][]byte{
			fakeTokenKey: []byte(fakeToken),
		},
	}

	It("Should create a ValidationResult and update its Status with a failed condition", func() {
		By("By creating a new MaasValidator")
		ctx := context.Background()

		Expect(k8sClient.Create(ctx, maasSecret)).Should(Succeed())
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
