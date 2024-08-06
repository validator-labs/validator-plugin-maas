// Package utils provides utility functions for the MAAS validator
package utils

import (
	"fmt"

	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vapiconstants "github.com/validator-labs/validator/pkg/constants"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
)

// BuildValidationResult builds a default ValidationResult for a given validation type
func BuildValidationResult(ruleName, ruleType string) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("Validation succeeded for %s", ruleType)
	latestCondition.ValidationRule = fmt.Sprintf("%s-%s", vapiconstants.ValidationRulePrefix, util.Sanitize(ruleName))
	latestCondition.ValidationType = ruleType
	return &types.ValidationRuleResult{Condition: &latestCondition, State: &state}
}

// UpdateResult updates a ValidationRuleResult with a list of errors and details
func UpdateResult(vr *types.ValidationRuleResult, errs []error, errMsg string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	vr.Condition.Details = append(vr.Condition.Details, details...)
}
