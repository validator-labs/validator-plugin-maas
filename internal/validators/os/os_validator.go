// Package os handles OS image rule reconciliation.
package os

import (
	"errors"
	"fmt"

	"github.com/canonical/gomaasclient/entity"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
)

// ImageAPI is an interface for interacting with the Maas API
type imageAPI interface {
	ListOSImages() ([]entity.BootResource, error)
}

// ImageRulesService is a service for reconciling OS image rules
type ImageRulesService struct {
	log logr.Logger
	api imageAPI
}

// NewImageRulesService returns a ImageRulesService
func NewImageRulesService(log logr.Logger, api imageAPI) *ImageRulesService {
	return &ImageRulesService{
		log: log,
		api: api,
	}
}

// ReconcileMaasInstanceImageRules reconciles a MaaS instance image rule from the MaasValidator config
func (s *ImageRulesService) ReconcileMaasInstanceImageRules(rule v1alpha1.OSImageRule) (*types.ValidationRuleResult, error) {

	vr := buildValidationResult(rule)

	brs, err := s.api.ListOSImages()
	if err != nil {
		return vr, err
	}

	errs, details := findBootResources(rule, brs)

	s.updateResult(vr, errs, constants.ErrMsg, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.OSImageRule) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.MaasInstance)
	latestCondition.ValidationRule = rule.Name
	latestCondition.ValidationType = constants.MaasInstance
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

// ListOSImages returns a list of OS images from the Maas API
func (s *ImageRulesService) ListOSImages() ([]entity.BootResource, error) {
	images, err := s.api.ListOSImages()
	if err != nil {
		return nil, err
	}
	return images, nil
}

// convertBootResourceToOSImage formats a list of BootResources as a list of OSImages
func convertBootResourceToOSImage(images []entity.BootResource) []v1alpha1.OSImage {
	converted := make([]v1alpha1.OSImage, len(images))
	for i, img := range images {
		converted[i] = v1alpha1.OSImage{
			OSName:       img.Name,
			Architecture: img.Architecture,
		}
	}
	return converted
}

// findBootResources checks if a list of OSImages is a subset of a list of BootResources
func findBootResources(rule v1alpha1.OSImageRule, images []entity.BootResource) (errs []error, details []string) {
	errs = make([]error, 0)
	details = make([]string, 0)

	converted := convertBootResourceToOSImage(images)
	convertedSet := mapset.NewSet(converted...)
	imgRulesSet := mapset.NewSet(rule.OSImages...)

	if imgRulesSet.IsSubset(convertedSet) {
		return errs, details
	}

	diffSet := imgRulesSet.Difference(convertedSet)

	diffSetIt := diffSet.Iterator()

	for img := range diffSetIt.C {
		errs = append(errs, errors.New(constants.ErrMsg))
		details = append(details, fmt.Sprintf("OS image %s with arch %s was not found", img.OSName, img.Architecture))
	}

	return errs, details
}
