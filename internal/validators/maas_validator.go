package validators

import (
	"fmt"

	"github.com/go-logr/logr"

	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util/ptr"
)

type MaasRuleService struct {
	log logr.Logger
}

func NewMaasRuleService(log logr.Logger) *MaasRuleService {
	return &MaasRuleService{
		log: log,
	}
}

// ReconcileMaasInstanceRule reconciles a MaaS instance rule from the MaasValidator config
func (s *MaasRuleService) ReconcileMaasInstanceRule(rule v1alpha1.MaasInstanceRule, username, password string) (*vapitypes.ValidationResult, error) {
	vr := buildValidationResult(rule)

	errs := make([]error, 0)
	details := make([]string, 0)

	errMsg := "failed to validate artifact in registry"
	s.updateResult(vr, errs, errMsg, rule.Name(), details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.MaasInstanceRule) *types.ValidationResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.MaasInstance)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = constants.MaasInstance
	return &types.ValidationResult{Condition: &latestCondition, State: &state}
}

func (s *MaasRuleService) updateResult(vr *types.ValidationResult, errs []error, errMsg, ruleName string, details ...string) {
	if len(errs) > 0 {
		vr.State = ptr.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	for _, detail := range details {
		vr.Condition.Details = append(vr.Condition.Details, detail)
	}
}
