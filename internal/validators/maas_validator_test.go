package validators

import (
	"errors"
	"testing"

	"github.com/maas/gomaasclient/entity"
	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/stretchr/testify/assert"
)

type DummyMaaSAPIClient struct {
	images      []entity.BootResource
	nameservers []v1alpha1.Nameserver
}

func (d *DummyMaaSAPIClient) ListDNSServers() ([]v1alpha1.Nameserver, error) {
	return d.nameservers, nil
}

func (d *DummyMaaSAPIClient) ListOSImages() ([]entity.BootResource, error) {

	return d.images, nil
}

func TestFindingBootResources(t *testing.T) {

	type TestCase struct {
		Name        string
		ruleService *MaasRuleService
		imageRules  []v1alpha1.OSImage
		errors      []error
		details     []string
	}

	testCases := []TestCase{
		{
			Name: "image is found in MaaS",
			ruleService: NewMaasRuleService(&DummyMaaSAPIClient{
				images: []entity.BootResource{
					{
						Name:         "Ubuntu",
						Architecture: "amd64/ga-20.04",
					},
				},
			}),
			imageRules: []v1alpha1.OSImage{
				{Name: "Ubuntu", Architecture: "amd64/ga-20.04"},
			},
			errors:  make([]error, 0),
			details: make([]string, 0),
		},
		{
			Name: "image is not found in MaaS",
			ruleService: NewMaasRuleService(
				&DummyMaaSAPIClient{
					images: make([]entity.BootResource, 0),
				}),
			imageRules: []v1alpha1.OSImage{
				{Name: "Ubuntu", Architecture: "amd64/ga-20.04"},
			},
			errors:  []error{errors.New("failed to validate rule")},
			details: []string{"OS image Ubuntu with arch amd64/ga-20.04 was not found"},
		},
		{
			Name: "a few images are not found in MaaS",
			ruleService: NewMaasRuleService(
				&DummyMaaSAPIClient{
					images: make([]entity.BootResource, 0),
				}),
			imageRules: []v1alpha1.OSImage{
				{Name: "Ubuntu", Architecture: "amd64/ga-20.04"},
				{Name: "Ubuntu", Architecture: "amd64/ga-22.04"},
			},
			errors: []error{errors.New("failed to validate rule"), errors.New("failed to validate rule")},
			details: []string{
				"OS image Ubuntu with arch amd64/ga-20.04 was not found",
				"OS image Ubuntu with arch amd64/ga-22.04 was not found"},
		},
	}

	for _, tc := range testCases {
		images, _ := tc.ruleService.ListOSImages()

		errors, details := findBootResources(tc.imageRules, images)

		assert.Equal(t, errors, tc.errors, tc.Name)
		for _, detail := range details {
			assert.Contains(t, tc.details, detail, tc.Name)
		}

	}

}

func TestExtDNS(t *testing.T) {

	type TestCase struct {
		Name             string
		ruleService      *MaasRuleService
		externalDNSRules []v1alpha1.Nameserver
		errors           []error
		details          []string
	}

	testCases := []TestCase{
		{Name: "MaaS has 2 external nameservers",
			ruleService: NewMaasRuleService(&DummyMaaSAPIClient{
				nameservers: []v1alpha1.Nameserver{"8.8.8.8", "9.9.9.9"},
			}),
			externalDNSRules: []v1alpha1.Nameserver{"8.8.8.8", "9.9.9.9"},
			errors:           make([]error, 0),
			details:          make([]string, 0),
		},
		{Name: "MaaS has 2 external nameservers, only 1 required",
			ruleService: NewMaasRuleService(&DummyMaaSAPIClient{
				nameservers: []v1alpha1.Nameserver{"8.8.8.8", "9.9.9.9"},
			}),
			externalDNSRules: []v1alpha1.Nameserver{"9.9.9.9"},
			errors:           make([]error, 0),
			details:          make([]string, 0),
		},
		{Name: "MaaS has 1 external nameservers, 2 required",
			ruleService: NewMaasRuleService(&DummyMaaSAPIClient{
				nameservers: []v1alpha1.Nameserver{"9.9.9.9"},
			}),
			externalDNSRules: []v1alpha1.Nameserver{"9.9.9.9", "8.8.8.8"},
			errors:           []error{errors.New("failed to validate rule")}, //, errors.New("failed to validate rule")},
			details:          []string{"External nameserver 8.8.8.8 was not found"},
		},
		{Name: "MaaS has 0 external nameservers, 2 required",
			ruleService: NewMaasRuleService(&DummyMaaSAPIClient{
				nameservers: make([]v1alpha1.Nameserver, 0),
			}),
			externalDNSRules: []v1alpha1.Nameserver{"9.9.9.9", "8.8.8.8"},
			errors:           []error{errors.New("failed to validate rule"), errors.New("failed to validate rule")},
			details:          []string{"External nameserver 8.8.8.8 was not found", "External nameserver 9.9.9.9 was not found"},
		},
		{Name: "MaaS has 2 external nameservers, 2 required",
			ruleService: NewMaasRuleService(&DummyMaaSAPIClient{
				nameservers: []v1alpha1.Nameserver{"1.1.1.1", "8.8.4.4"},
			}),
			externalDNSRules: []v1alpha1.Nameserver{"9.9.9.9", "8.8.8.8"},
			errors:           []error{errors.New("failed to validate rule"), errors.New("failed to validate rule")},
			details:          []string{"External nameserver 8.8.8.8 was not found", "External nameserver 9.9.9.9 was not found"},
		},
	}
	for _, tc := range testCases {
		nameservers, _ := tc.ruleService.apiclient.ListDNSServers()

		errors, details := assertExternalDNS(tc.externalDNSRules, nameservers)
		assert.Equal(t, errors, tc.errors, tc.Name)
		for _, detail := range details {
			assert.Contains(t, tc.details, detail, tc.Name)
		}

	}
}
