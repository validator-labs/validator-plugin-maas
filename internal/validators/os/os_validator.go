// Package os handles OS image rule reconciliation.
package os

import (
	"fmt"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vapiconstants "github.com/validator-labs/validator/pkg/constants"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
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

	vr := buildValidationResult(rule)

	errs, details := s.findBootResources(rule)

	s.updateResult(vr, errs, constants.ErrImageNotFound, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.ImageRule) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.ValidationTypeImage)
	latestCondition.ValidationRule = fmt.Sprintf("%s-%s", vapiconstants.ValidationRulePrefix, util.Sanitize(rule.Name))
	latestCondition.ValidationType = constants.ValidationTypeImage
	return &types.ValidationRuleResult{Condition: &latestCondition, State: &state}
}

// updateResult updates a ValidationRuleResult with a list of errors and details
func (s *ImageRulesService) updateResult(vr *types.ValidationRuleResult, errs []error, errMsg string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	vr.Condition.Details = append(vr.Condition.Details, details...)
}

// convertBootResourceToOSImage formats a list of BootResources as a list of OSImages
func convertBootResourceToOSImage(images []entity.BootResource) []v1alpha1.Image {
	converted := make([]v1alpha1.Image, len(images))
	for i, img := range images {
		converted[i] = v1alpha1.Image{
			Name:         img.Name,
			Architecture: img.Architecture,
		}
	}
	return converted
}

// findBootResources checks if a list of Images is a subset of a list of BootResources
func (s *ImageRulesService) findBootResources(rule v1alpha1.ImageRule) (errs []error, details []string) {
	errs = make([]error, 0)
	details = make([]string, 0)

	images, err := s.api.Get(&entity.BootResourcesReadParams{})
	if err != nil {
		return errs, details
	}
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
	return errs, details
}
