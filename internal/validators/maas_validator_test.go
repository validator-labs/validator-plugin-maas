package validators

import (
	"testing"

	"github.com/maas/gomaasclient/entity"
	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/stretchr/testify/assert"
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
		ruleService *MaasRuleService
		imageRule   v1alpha1.OSImage
		errors      []error
		details     []string
	}

	tc := TestCase{
		ruleService: NewMaasRuleService(&DummyMaaSAPIClient{
			images: []entity.BootResource{
				{
					Name:         "Ubuntu",
					Architecture: "amd64/ga-20.04",
				},
			},
		}),
		imageRule: v1alpha1.OSImage{
			Name:         "Ubuntu",
			Architecture: "amd64/ga-20.04",
		},
		errors:  make([]error, 0),
		details: make([]string, 0),
	}

	images, _ := tc.ruleService.ListOSImages()
	assert.Equal(t, len(images), 1)

	errors, details := findBootResources(tc.imageRule, images)

	assert.Equal(t, errors, tc.errors)
	assert.Equal(t, details, tc.details)
}
