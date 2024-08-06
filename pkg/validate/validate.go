// Package validate defines a Validate function that evaluates a MaasValidatorSpec and returns a ValidationResponse.
package validate

import (
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"

	dnsval "github.com/validator-labs/validator-plugin-maas/pkg/validators/dns"
	osval "github.com/validator-labs/validator-plugin-maas/pkg/validators/os"
	resval "github.com/validator-labs/validator-plugin-maas/pkg/validators/res"

	maasclient "github.com/canonical/gomaasclient/client"
)

// SetUpClient is defined to enable monkey patching the setUpClient function in integration tests
var SetUpClient = setUpClient

// Validate validates the MaasValidatorSpec and returns a ValidationResponse.
func Validate(spec v1alpha1.MaasValidatorSpec, maasURL string, maasToken string, log logr.Logger) types.ValidationResponse {

	maasClient, err := SetUpClient(maasURL, maasToken)
	if err != nil {
		log.Error(err, "failed to initialize MAAS client")
	}

	resp := types.ValidationResponse{}

	imageRulesService := osval.NewImageRulesService(log, maasClient.BootResources)
	resourceRulesService := resval.NewResourceRulesService(log, maasClient.Machines)
	upstreamDNSRulesService := dnsval.NewUpstreamDNSRulesService(log, maasClient.MAASServer)
	internalDNSRulesService := dnsval.NewInternalDNSRulesService(log, maasClient.DNSResources)

	// MAAS Instance image rules
	for _, rule := range spec.ImageRules {
		vrr, err := imageRulesService.ReconcileMaasInstanceImageRule(rule)
		if err != nil {
			log.V(0).Error(err, "failed to reconcile MAAS image rule")
		}
		resp.AddResult(vrr, err)
	}

	// MAAS Instance upstream DNS rules
	for _, rule := range spec.UpstreamDNSRules {
		vrr, err := upstreamDNSRulesService.ReconcileMaasInstanceUpstreamDNSRule(rule)
		if err != nil {
			log.V(0).Error(err, "failed to reconcile MAAS upstream DNS rule")
		}
		resp.AddResult(vrr, err)
	}

	seenAZ := make(map[string]bool, 0)
	// MAAS Instance resource availability rules
	for _, rule := range spec.ResourceAvailabilityRules {
		vrr, err := resourceRulesService.ReconcileMaasInstanceResourceRule(rule, seenAZ)
		if err != nil {
			log.V(0).Error(err, "failed to reconcile MAAS resource rule")
		}
		resp.AddResult(vrr, err)
		seenAZ[rule.AZ] = true
	}

	// MAAS Instance internal DNS rules
	for _, rule := range spec.InternalDNSRules {
		vrr, err := internalDNSRulesService.ReconcileMaasInstanceInternalDNSRule(rule)
		if err != nil {
			log.V(0).Error(err, "failed to reconcile MAAS internal DNS rule")
		}
		resp.AddResult(vrr, err)
	}

	return resp
}

func setUpClient(maasURL, maasToken string) (*maasclient.Client, error) {
	maasClient, err := maasclient.GetClient(maasURL, maasToken, "2.0")
	if err != nil {
		return nil, err
	}
	return maasClient, nil
}
