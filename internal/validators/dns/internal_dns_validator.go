package dns

import (
	"fmt"
	"slices"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	"github.com/validator-labs/validator-plugin-maas/internal/utils"
	"github.com/validator-labs/validator/pkg/types"
)

// InternalDNSRulesService is the service for validating MAAS instance internal DNS rules
type InternalDNSRulesService struct {
	log logr.Logger
	api api.DNSResources
}

// NewInternalDNSRulesService creates a new InternalDNSRulesService
func NewInternalDNSRulesService(log logr.Logger, api api.DNSResources) *InternalDNSRulesService {
	return &InternalDNSRulesService{
		log: log,
		api: api,
	}
}

// ReconcileMaasInstanceInternalDNSRule reconciles a MAAS instance internal DNS rule
func (s *InternalDNSRulesService) ReconcileMaasInstanceInternalDNSRule(rule v1alpha1.InternalDNSRule) (*types.ValidationRuleResult, error) {
	vr := utils.BuildValidationResult(rule.MaasDomain, constants.ValidationTypeIDNS)

	details, errs := s.findDNSServers(rule)

	utils.UpdateResult(vr, errs, constants.ErrIDNSNotConfigured, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}

	return vr, nil
}

func (s *InternalDNSRulesService) findDNSServers(rule v1alpha1.InternalDNSRule) ([]string, []error) {
	details := make([]string, 0)
	errs := make([]error, 0)

	maasDNSResources, err := s.api.Get(&entity.DNSResourcesGetParams{All: true})
	if err != nil {
		return nil, []error{err}
	}

	formattedRecords := formatDNSRecords(maasDNSResources)

	for _, res := range rule.DNSResources {
		if !checkOneResourcePresent(formattedRecords, res) {
			errs = append(errs, fmt.Errorf("one or more DNS records not found for %s", res.FQDN))
		} else {
			details = append(details, fmt.Sprintf("All required DNS records found for %s", res.FQDN))
		}
	}

	return details, errs
}

func checkOneResourcePresent(formattedRecords map[string][]v1alpha1.DNSRecord, resource v1alpha1.DNSResource) bool {
	if formattedRecords[resource.FQDN] == nil {
		return false
	}

	for _, rec := range resource.DNSRecords {
		frr := formattedRecords[resource.FQDN]
		if !slices.Contains(frr, rec) {
			return false
		}
	}

	return true
}

func formatDNSRecords(dnsResources []entity.DNSResource) map[string][]v1alpha1.DNSRecord {
	formattedRecords := make(map[string][]v1alpha1.DNSRecord)

	for _, r := range dnsResources {
		if r.ResourceRecords != nil && len(r.ResourceRecords) > 0 {
			fr := make([]v1alpha1.DNSRecord, 0)
			for _, rr := range r.ResourceRecords {
				fr = append(fr, v1alpha1.DNSRecord{
					IP:   rr.RRData,
					Type: rr.RRType,
					TTL:  rr.TTL,
				})
			}
			formattedRecords[r.FQDN] = fr
		}
	}
	return formattedRecords
}
