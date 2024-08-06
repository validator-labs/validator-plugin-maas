package res

import (
	"testing"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
)

type DummyMachine struct {
	api.Machines
	MachinesList []entity.Machine
}

func (d *DummyMachine) Get(params *entity.MachinesParams) ([]entity.Machine, error) {
	return d.MachinesList, nil
}

func TestReconcileMaasInstanceImageRule(t *testing.T) {

	testCases := []struct {
		Name        string
		ruleService *ResourceRulesService
		resources   []v1alpha1.ResourceAvailabilityRule
		errors      []string
		details     []string
	}{
		{
			Name: "Enough machines are found in MAAS",
			ruleService: NewResourceRulesService(
				logr.Logger{},
				&DummyMachine{
					MachinesList: []entity.Machine{
						{Zone: entity.Zone{Name: "az1"}, Hostname: "maas.foo", ResourceURI: "/api/2.0/machines/1/", CPUCount: 24, Memory: 32 * 1024, Storage: 150 * 1000, Pool: entity.ResourcePool{Name: "pool1"}, TagNames: []string{"tag1", "tag2"}},
						{Zone: entity.Zone{Name: "az1"}, Hostname: "maas.foo", ResourceURI: "/api/2.0/machines/1/", CPUCount: 24, Memory: 32 * 1024, Storage: 150 * 1000, Pool: entity.ResourcePool{Name: "pool1"}, TagNames: []string{"tag1", "tag2"}},
						{Zone: entity.Zone{Name: "az1"}, Hostname: "maas.foo", ResourceURI: "/api/2.0/machines/1/", CPUCount: 24, Memory: 32 * 1024, Storage: 150 * 1000, Pool: entity.ResourcePool{Name: "pool1"}, TagNames: []string{"tag1", "tag2"}},
						{Zone: entity.Zone{Name: "az1"}, Hostname: "maas.foo", ResourceURI: "/api/2.0/machines/1/", CPUCount: 24, Memory: 32 * 1024, Storage: 150 * 1000, Pool: entity.ResourcePool{Name: "pool1"}, TagNames: []string{"tag1", "tag2"}},
					},
				},
			),
			resources: []v1alpha1.ResourceAvailabilityRule{
				{Name: "AZ1 rule 1", AZ: "az1", Resources: []v1alpha1.Resource{
					{NumMachines: 3, NumCPU: 16, RAM: 16, Disk: 100, Pool: "pool1", Tags: []string{"tag1", "tag2"}},
				},
				},
			},
			errors:  nil,
			details: []string{"Found 3 machine(s) available with 16 Cores, 16GB RAM, 100GB Disk"},
		},
		{
			Name: "Not enough machines are found in MAAS",
			ruleService: NewResourceRulesService(
				logr.Logger{},
				&DummyMachine{
					MachinesList: []entity.Machine{
						{Zone: entity.Zone{Name: "az1"}, Hostname: "maas.foo", ResourceURI: "/api/2.0/machines/1/", CPUCount: 24, Memory: 32 * 1024, Storage: 150 * 1000, Pool: entity.ResourcePool{Name: "pool1"}, TagNames: []string{"tag1", "tag2"}},
						{Zone: entity.Zone{Name: "az1"}, Hostname: "maas.foo", ResourceURI: "/api/2.0/machines/1/", CPUCount: 12, Memory: 32 * 1024, Storage: 150 * 1000, Pool: entity.ResourcePool{Name: "pool1"}, TagNames: []string{"tag1", "tag2"}},
						{Zone: entity.Zone{Name: "az1"}, Hostname: "maas.foo", ResourceURI: "/api/2.0/machines/1/", CPUCount: 12, Memory: 32 * 1024, Storage: 150 * 1000, Pool: entity.ResourcePool{Name: "pool1"}, TagNames: []string{"tag1", "tag2"}},
					},
				}),
			resources: []v1alpha1.ResourceAvailabilityRule{
				{Name: "AZ1 rule 2", AZ: "az1", Resources: []v1alpha1.Resource{
					{NumMachines: 2, NumCPU: 16, RAM: 16, Disk: 100, Pool: "pool1", Tags: []string{"tag1", "tag2"}},
				}},
			},
			errors:  []string{"insufficient machines available with 16 Cores, 16GB RAM, 100GB Disk. 1/2 available"},
			details: nil,
		},
		{
			Name: "No machines are found in MAAS",
			ruleService: NewResourceRulesService(
				logr.Logger{},
				&DummyMachine{
					MachinesList: []entity.Machine{},
				}),
			resources: []v1alpha1.ResourceAvailabilityRule{
				{Name: "AZ1 rule 2", AZ: "az1", Resources: []v1alpha1.Resource{
					{NumMachines: 1, NumCPU: 16, RAM: 16, Disk: 100, Pool: "pool1", Tags: []string{"tag1", "tag2"}},
				}},
			},
			errors:  []string{"not enough resources available in az: have: 0, need: 1"},
			details: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var details []string
			var errors []string

			for _, rule := range tc.resources {
				vr, _ := tc.ruleService.ReconcileMaasInstanceResourceRule(rule, map[string]bool{})
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
