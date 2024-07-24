// Package dns contains the logic for validating MAAS instance DNS rules
package dns

import (
	"fmt"
	"strings"

	"github.com/canonical/gomaasclient/api"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	"github.com/validator-labs/validator-plugin-maas/internal/utils"
	"github.com/validator-labs/validator/pkg/types"
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

	vr := utils.BuildValidationResult(rule.Name, constants.ValidationTypeUDNS)

	details, errs := s.findDNSServers(rule.NumDNSServers)

	utils.UpdateResult(vr, errs, constants.ErrUDNSNotConfigured, details...)

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
