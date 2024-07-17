package os

import (
	"errors"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
)

type DummyMaaSAPIClient struct {
	images []entity.BootResource
}

func (d *DummyMaaSAPIClient) ListOSImages() ([]entity.BootResource, error) {
	return d.images, nil
}

func TestFindingBootResources(t *testing.T) {

	type TestCase struct {
		Name        string
		ruleService *ImageRulesService
		imageRules  []v1alpha1.OSImageRule
		errors      []error
		details     []string
	}

	testCases := []TestCase{
		{
			Name: "image is found in MaaS",
			ruleService: NewImageRulesService(
				logr.Logger{},
				&DummyMaaSAPIClient{
					images: []entity.BootResource{
						{
							Name:         "Ubuntu",
							Architecture: "amd64/ga-20.04",
						},
					},
				},
			),
			imageRules: []v1alpha1.OSImageRule{
				{Name: "Image found", OSImages: []v1alpha1.OSImage{
					{OSName: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
			},
			errors:  make([]error, 0),
			details: make([]string, 0),
		},
		{
			Name: "image is not found in MaaS",
			ruleService: NewImageRulesService(
				logr.Logger{},
				&DummyMaaSAPIClient{
					images: make([]entity.BootResource, 0),
				}),
			imageRules: []v1alpha1.OSImageRule{
				{Name: "Image not found", OSImages: []v1alpha1.OSImage{
					{OSName: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
			},
			errors:  []error{errors.New("failed to validate rule")},
			details: []string{"OS image Ubuntu with arch amd64/ga-20.04 was not found"},
		},
		{
			Name: "a few images are not found in MaaS",
			ruleService: NewImageRulesService(
				logr.Logger{},
				&DummyMaaSAPIClient{
					images: make([]entity.BootResource, 0),
				}),
			imageRules: []v1alpha1.OSImageRule{
				{Name: "Image not found", OSImages: []v1alpha1.OSImage{
					{OSName: "Ubuntu", Architecture: "amd64/ga-20.04"},
				}},
				{Name: "Image not found", OSImages: []v1alpha1.OSImage{
					{OSName: "Ubuntu", Architecture: "amd64/ga-22.04"},
				}},
			},
			errors: []error{errors.New("failed to validate rule"), errors.New("failed to validate rule")},
			details: []string{
				"OS image Ubuntu with arch amd64/ga-20.04 was not found",
				"OS image Ubuntu with arch amd64/ga-22.04 was not found"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var allErrors []error
			var allDetails []string

			for _, rule := range tc.imageRules {
				images, _ := tc.ruleService.ListOSImages()
				errors, details := findBootResources(rule, images)
				allErrors = append(allErrors, errors...)
				allDetails = append(allDetails, details...)
			}

			assert.Equal(t, len(tc.errors), len(allErrors), "Number of errors should match")
			for _, expectedError := range tc.errors {
				assert.Contains(t, allErrors, expectedError, "Expected error should be present")
			}

			assert.Equal(t, len(tc.details), len(allDetails), "Number of details should match")
			for _, expectedDetail := range tc.details {
				assert.Contains(t, allDetails, expectedDetail, "Expected detail should be present")
			}
		})
	}
}
