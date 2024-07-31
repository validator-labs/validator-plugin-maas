package os

import (
	"testing"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
)

type DummyBootResources struct {
	api.BootResources
	images []entity.BootResource
}

func (d *DummyBootResources) Get(params *entity.BootResourcesReadParams) ([]entity.BootResource, error) {
	return d.images, nil
}

func TestReconcileMaasInstanceImageRule(t *testing.T) {

	testCases := []struct {
		Name        string
		ruleService *ImageRulesService
		imageRules  []v1alpha1.ImageRule
		errors      []string
		details     []string
	}{
		{
			Name: "image is found in MAAS",
			ruleService: NewImageRulesService(
				logr.Logger{},
				&DummyBootResources{
					images: []entity.BootResource{
						{
							Type:         "Synced",
							Name:         "Ubuntu",
							Architecture: "amd64/ga-20.04",
						},
					},
				},
			),
			imageRules: []v1alpha1.ImageRule{
				{Name: "Image found", Images: []v1alpha1.Image{
					{Name: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
			},
			errors:  nil,
			details: []string{"OS image Ubuntu with arch amd64/ga-20.04 was found"},
		},
		{
			Name: "image is not found in MAAS",
			ruleService: NewImageRulesService(
				logr.Logger{},
				&DummyBootResources{
					images: make([]entity.BootResource, 0),
				}),
			imageRules: []v1alpha1.ImageRule{
				{Name: "Image not found", Images: []v1alpha1.Image{
					{Name: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
			},
			errors:  []string{"OS image Ubuntu with arch amd64/ga-20.04 was not found"},
			details: nil,
		},
		{
			Name: "a few images are not found in MAAS",
			ruleService: NewImageRulesService(
				logr.Logger{},
				&DummyBootResources{
					images: make([]entity.BootResource, 0),
				}),
			imageRules: []v1alpha1.ImageRule{
				{Name: "Image not found", Images: []v1alpha1.Image{
					{Name: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
				{Name: "Image not found", Images: []v1alpha1.Image{
					{Name: "Ubuntu", Architecture: "amd64/ga-22.04"},
				}},
			},
			errors:  []string{"OS image Ubuntu with arch amd64/ga-20.04 was not found", "OS image Ubuntu with arch amd64/ga-22.04 was not found"},
			details: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var details []string
			var errors []string

			for _, rule := range tc.imageRules {
				vr, _ := tc.ruleService.ReconcileMaasInstanceImageRule(rule)
				details = append(details, vr.Condition.Details...)
				errors = append(errors, vr.Condition.Failures...)
			}

			assert.Equal(t, len(tc.errors), len(errors), "Number of errors should match")
			for _, expectedError := range tc.errors {
				assert.Contains(t, errors, expectedError, "Expected error should be present")
			}
			assert.Equal(t, len(tc.details), len(details), "Number of details should match")
			for _, expectedDetail := range tc.details {
				assert.Contains(t, details, expectedDetail, "Expected detail should be present")
			}
		})
	}
}
