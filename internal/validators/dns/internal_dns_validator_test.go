package dns

import (
	"testing"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
)

type DummyDNSResources struct {
	api.DNSResources
	resources []entity.DNSResource
}

func (d *DummyDNSResources) Get(params *entity.DNSResourcesGetParams) ([]entity.DNSResource, error) {
	return d.resources, nil
}

func TestReconcileMaasInternalDNSRule(t *testing.T) {

	testCases := []struct {
		Name             string
		ruleService      *InternalDNSRulesService
		internalDNSRules []v1alpha1.InternalDNSRule
		errors           []string
		details          []string
	}{
		{
			Name: "Enough DNS servers are found in MAAS",
			ruleService: NewInternalDNSRulesService(
				logr.Logger{},
				&DummyDNSResources{
					resources: []entity.DNSResource{
						{
							FQDN: "foo.test.com", ResourceRecords: []entity.DNSResourceRecord{
								{
									TTL:    10,
									RRData: "1.1.1.1",
									RRType: "A",
								},
								{
									TTL:    10,
									RRData: "1.1.1.2",
									RRType: "A",
								},
							},
						},
						{
							FQDN: "bar.test.com", ResourceRecords: []entity.DNSResourceRecord{
								{
									TTL:    10,
									RRData: "1.1.1.1",
									RRType: "A",
								},
							},
						},
					},
				},
			),
			internalDNSRules: []v1alpha1.InternalDNSRule{
				{MaasDomain: "test.com", DNSResources: []v1alpha1.DNSResource{
					{FQDN: "foo.test.com", DNSRecords: []v1alpha1.DNSRecord{
						{
							TTL:  10,
							IP:   "1.1.1.1",
							Type: "A",
						},
					}},
					{FQDN: "bar.test.com", DNSRecords: []v1alpha1.DNSRecord{
						{
							TTL:  10,
							IP:   "1.1.1.1",
							Type: "A",
						},
					}},
				}},
			},
			errors:  nil,
			details: []string{"All required DNS records found for foo.test.com", "All required DNS records found for bar.test.com"},
		},
		{
			Name: "Not enough DNS servers are found in MAAS",
			ruleService: NewInternalDNSRulesService(
				logr.Logger{},
				&DummyDNSResources{
					resources: []entity.DNSResource{
						{
							FQDN: "foo.test.com", ResourceRecords: []entity.DNSResourceRecord{
								{
									TTL:    10,
									RRData: "1.1.1.1",
									RRType: "A",
								},
							},
						},
					},
				}),
			internalDNSRules: []v1alpha1.InternalDNSRule{
				{MaasDomain: "test.com", DNSResources: []v1alpha1.DNSResource{
					{FQDN: "foo.test.com", DNSRecords: []v1alpha1.DNSRecord{
						{
							TTL:  10,
							IP:   "1.1.1.2",
							Type: "A",
						},
					}},
				}},
			},
			errors:  []string{"one or more DNS records not found for foo.test.com"},
			details: nil,
		},
		{
			Name: "No DNS servers are found in MAAS",
			ruleService: NewInternalDNSRulesService(
				logr.Logger{},
				&DummyDNSResources{
					resources: []entity.DNSResource{},
				}),
			internalDNSRules: []v1alpha1.InternalDNSRule{
				{MaasDomain: "test.com", DNSResources: []v1alpha1.DNSResource{
					{FQDN: "foo.test.com", DNSRecords: []v1alpha1.DNSRecord{
						{
							TTL:  10,
							IP:   "1.1.1.2",
							Type: "A",
						},
					}},
				}},
			},
			errors:  []string{"one or more DNS records not found for foo.test.com"},
			details: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var details []string
			var errors []string

			for _, rule := range tc.internalDNSRules {
				vr, _ := tc.ruleService.ReconcileMaasInstanceInternalDNSRule(rule)
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
