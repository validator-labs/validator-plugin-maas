// Package dns contains the logic for validating MAAS instance DNS rules
package dns

import (
	"fmt"
	"strings"

	"github.com/canonical/gomaasclient/api"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vapiconstants "github.com/validator-labs/validator/pkg/constants"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
)

// UpstreamDNSRulesService is the service for validating MAAS instance upstream DNS rules
type UpstreamDNSRulesService struct {
	log logr.Logger
	api api.MAASServer
}

// NewUpstreamDNSRulesService creates a new UpstreamDNSRulesService
func NewUpstreamDNSRulesService(log logr.Logger, api api.MAASServer) *UpstreamDNSRulesService {
	return &UpstreamDNSRulesService{
		log: log,
		api: api,
	}
}

// ReconcileMaasInstanceUpstreamDNSRules reconciles a MAAS instance upstream DNS rule
func (s *UpstreamDNSRulesService) ReconcileMaasInstanceUpstreamDNSRules(rule v1alpha1.UpstreamDNSRule) (*types.ValidationRuleResult, error) {

	vr := buildValidationResult(rule)

	details, errs := s.findDNSServers(rule.NumDNSServers)

	updateResult(vr, errs, constants.ErrUDNSNotConfigured, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}

	return vr, nil
}

func (s *UpstreamDNSRulesService) findDNSServers(expected int) ([]string, []error) {
	details := make([]string, 0)
	errs := make([]error, 0)

	ns, err := s.api.Get("upstream_dns")
	if err != nil {
		return nil, []error{err}
	}
	nameservers := strings.Split(string(ns), " ")
	numServers := len(nameservers)

	if nameservers[0] == "" {
		numServers = 0
	}

	if numServers < expected {
		errs = append(errs, fmt.Errorf("expected %d DNS server(s), got %d", expected, numServers))
	} else {
		details = append(details, fmt.Sprintf("Found %d DNS server(s)", len(nameservers)))
	}
	return details, errs
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.UpstreamDNSRule) *types.ValidationRuleResult {
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
func updateResult(vr *types.ValidationRuleResult, errs []error, errMsg string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	vr.Condition.Details = append(vr.Condition.Details, details...)
}
