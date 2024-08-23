package dns

import (
	"testing"

	"github.com/canonical/gomaasclient/api"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
)

type DummyMAASServer struct {
	api.MAASServer
	upstreamDNS string
}

func (d *DummyMAASServer) Get(string) ([]byte, error) {
	return []byte(d.upstreamDNS), nil
}

func TestReconcileMaasInstanceUpstreamDNSRule(t *testing.T) {

	testCases := []struct {
		Name             string
		ruleService      *UpstreamDNSRulesService
		upstreamDNSRules []v1alpha1.UpstreamDNSRule
		errors           []string
		details          []string
	}{
		{
			Name: "Enough DNS servers are found in MAAS",
			ruleService: NewUpstreamDNSRulesService(
				logr.Logger{},
				&DummyMAASServer{
					upstreamDNS: "8.8.8.8",
				},
			),
			upstreamDNSRules: []v1alpha1.UpstreamDNSRule{
				{RuleName: "Upstream DNS rule 1", NumDNSServers: 1},
			},
			errors:  nil,
			details: []string{"Found 1 DNS server(s)"},
		},
		{
			Name: "Not enough DNS servers are found in MAAS",
			ruleService: NewUpstreamDNSRulesService(
				logr.Logger{},
				&DummyMAASServer{
					upstreamDNS: "8.8.8.8",
				}),
			upstreamDNSRules: []v1alpha1.UpstreamDNSRule{
				{RuleName: "Upstream DNS rule 2", NumDNSServers: 2},
			},
			errors:  []string{"expected 2 DNS server(s), got 1"},
			details: nil,
		},
		{
			Name: "No DNS servers are found in MAAS",
			ruleService: NewUpstreamDNSRulesService(
				logr.Logger{},
				&DummyMAASServer{
					upstreamDNS: "",
				}),
			upstreamDNSRules: []v1alpha1.UpstreamDNSRule{
				{RuleName: "Upstream DNS rule 3", NumDNSServers: 1},
			},
			errors:  []string{"expected 1 DNS server(s), got 0"},
			details: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var details []string
			var errors []string

			for _, rule := range tc.upstreamDNSRules {
				vr, _ := tc.ruleService.ReconcileMaasInstanceUpstreamDNSRule(rule)
				details = append(details, vr.Condition.Details...)
				errors = append(errors, vr.Condition.Failures...)
			}

			assert.Equal(t, len(tc.errors), len(errors), "Number of errors should match")
			for _, expectedError := range tc.errors {
				assert.Contains(t, errors, expectedError, "Expected error should be present")
			}
			assert.Equal(t, len(tc.details), len(details), "Number of details should match")
			for _, expectedDetail := range tc.details {
				assert.Contains(t, details, expectedDetail, "Expected detail should be present")
			}
		})
	}
}
