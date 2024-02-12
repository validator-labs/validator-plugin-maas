package validators

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/maas/gomaasclient/entity"

	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util"
)

type MaasRuleService struct {
	imagereader OSImageReader
	log         logr.Logger
}

type OSImageReader interface {
	Get(params *entity.BootResourcesReadParams) ([]entity.BootResource, error)
}

func NewMaasRuleService(log logr.Logger, imagereader OSImageReader) *MaasRuleService {
	return &MaasRuleService{
		imagereader: imagereader,
		log:         log,
	}
}

// ReconcileMaasInstanceRule reconciles a MaaS instance rule from the MaasValidator config
func (s *MaasRuleService) ReconcileMaasInstanceRule(imgRule v1alpha1.OSImage) (*vapitypes.ValidationResult, error) {
	vr := buildValidationResult(imgRule)

	errMsg := "failed to validate rule"
	errs := make([]error, 0)
	details := make([]string, 0)

	brs, err := s.listOSImages()
	if err != nil {
		return vr, err
	}
	var found bool = false
	for _, br := range brs {
		if (imgRule.Name == br.Name) && (imgRule.Architecture == br.Architecture) {
			found = true
			break
		}
	}
	if !found {
		errs = append(errs, errors.New(errMsg))
		details = append(details, fmt.Sprintf("OS image %s with arch %s was not found", imgRule.Name, imgRule.Architecture))
	}
	s.updateResult(vr, errs, errMsg, imgRule.Name, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.OSImage) *types.ValidationResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.MaasInstance)
	latestCondition.ValidationRule = rule.Name
	latestCondition.ValidationType = constants.MaasInstance
	return &types.ValidationResult{Condition: &latestCondition, State: &state}
}

func (s *MaasRuleService) updateResult(vr *types.ValidationResult, errs []error, errMsg, ruleName string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	for _, detail := range details {
		vr.Condition.Details = append(vr.Condition.Details, detail)
	}
}

func (s *MaasRuleService) listOSImages() ([]entity.BootResource, error) {
	bootResoursces := make([]entity.BootResource, 0)
	readEntity := entity.BootResourcesReadParams{}
	br, err := s.imagereader.Get(&readEntity)

	if err != nil {
		return bootResoursces, err
	}
	return br, nil
	//for _, b := range br {
	//	fmt.Println(b.Architecture, b.Name)
	//}
}
