// Package os handles OS image rule reconciliation.
package os

import (
	"fmt"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/utils"
	"github.com/validator-labs/validator-plugin-maas/pkg/constants"
	"github.com/validator-labs/validator/pkg/types"
)

// ImageRulesService is a service for reconciling OS image rules
type ImageRulesService struct {
	log logr.Logger
	api api.BootResources
}

// NewImageRulesService returns a ImageRulesService
func NewImageRulesService(log logr.Logger, api api.BootResources) *ImageRulesService {
	return &ImageRulesService{
		log: log,
		api: api,
	}
}

// ReconcileMaasInstanceImageRule reconciles a MAAS instance image rule from the MaasValidator config
func (s *ImageRulesService) ReconcileMaasInstanceImageRule(rule v1alpha1.ImageRule) (*types.ValidationRuleResult, error) {

	vr := utils.BuildValidationResult(rule.Name)

	details, errs := s.findBootResources(rule)

	utils.UpdateResult(vr, errs, constants.ErrImageNotFound, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// convertBootResourceToOSImage formats a list of BootResources as a list of OSImages
func convertBootResourceToOSImage(images []entity.BootResource) []v1alpha1.Image {
	converted := make([]v1alpha1.Image, 0)
	for _, img := range images {
		// the client lib does not seem to care about params it is given, so for now manually filtering
		if img.Type == "Synced" {
			converted = append(converted, v1alpha1.Image{
				Name:         img.Name,
				Architecture: img.Architecture,
			})
		}
	}
	return converted
}

// findBootResources checks if a list of Images is a subset of a list of BootResources
func (s *ImageRulesService) findBootResources(rule v1alpha1.ImageRule) ([]string, []error) {
	images, err := s.api.Get(&entity.BootResourcesReadParams{Type: "Synced"})
	if err != nil {
		return nil, []error{err}
	}

	details := make([]string, 0)
	errs := make([]error, 0)

	converted := convertBootResourceToOSImage(images)
	convertedSet := mapset.NewSet(converted...)
	imgRulesSet := mapset.NewSet(rule.Images...)

	diffSet := imgRulesSet.Difference(convertedSet)
	intersectionSet := imgRulesSet.Intersect(convertedSet)

	for img := range diffSet.Iterator().C {
		errs = append(errs, fmt.Errorf("OS image %s with arch %s was not found", img.Name, img.Architecture))
	}

	for img := range intersectionSet.Iterator().C {
		details = append(details, fmt.Sprintf("OS image %s with arch %s was found", img.Name, img.Architecture))
	}
	return details, errs
}
