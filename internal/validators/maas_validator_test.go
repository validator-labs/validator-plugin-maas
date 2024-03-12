package validators

import (
	"errors"
	"testing"

	"github.com/maas/gomaasclient/entity"
	"github.com/stretchr/testify/assert"

	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
)

type DummyMaaSAPIClient struct {
	images []entity.BootResource
}

func (d *DummyMaaSAPIClient) ListDNSServers() ([]entity.DNSResource, error) {
	return make([]entity.DNSResource, 0), nil
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
