// Package res contains the logic for validating resource rules
package res

import (
	"fmt"
	"sort"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	"github.com/validator-labs/validator-plugin-maas/internal/utils"
	"github.com/validator-labs/validator/pkg/types"
)

type resource struct {
	NumCPU int
	RAM    int
	Disk   int
	Pool   string
	Tags   []string
}

// ResourceRulesService is a service for reconciling resource rules
type ResourceRulesService struct {
	log logr.Logger
	api api.Machines
}

// NewResourceRulesService returns a ResourceRulesService
func NewResourceRulesService(log logr.Logger, api api.Machines) *ResourceRulesService {
	return &ResourceRulesService{
		log: log,
		api: api,
	}
}

// ReconcileMaasInstanceResourceRule reconciles a MAAS instance resource rule from the MaasValidator config
func (s *ResourceRulesService) ReconcileMaasInstanceResourceRule(rule v1alpha1.ResourceAvailabilityRule, seen []string) (*types.ValidationRuleResult, error) {
	errs := make([]error, 0)
	details := make([]string, 0)

	vr := utils.BuildValidationResult(rule.Name, constants.ValidationTypeResource)

	// do not process an AZ more than once
	if containsString(seen, rule.AZ) {
		errs = append(errs, fmt.Errorf("availability zone %s already validated", rule.AZ))
	} else {
		errs, details = s.findMachineResources(rule)
	}

	utils.UpdateResult(vr, errs, constants.ErrImageNotFound, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// contains checks if a slice contains a string
func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// findMachineResources checks whether the rule is satisfied by the machines in the availability zone
func (s *ResourceRulesService) findMachineResources(rule v1alpha1.ResourceAvailabilityRule) ([]error, []string) {
	errs := make([]error, 0)
	details := make([]string, 0)

	// Get all "ready" machines in the availability zone
	machines, err := s.api.Get(&entity.MachinesParams{Zone: []string{rule.AZ}, Status: []string{"ready"}})
	if err != nil {
		errs = append(errs, fmt.Errorf("error retrieving machines in availability zone %s", rule.AZ))
		return errs, details
	}

	resources := formatMachines(machines)

	// each rule may have multiple resource checks for 1 AZ
	for _, rr := range rule.Resources {
		detail, err := compareResources(&resources, &rr)
		if err != nil {
			errs = append(errs, err)
		}
		if detail != "" {
			details = append(details, detail)
		}
	}

	return errs, details
}

func compareResources(resources *[]resource, expected *v1alpha1.Resource) (string, error) {
	need := expected.NumMachines
	errMsg := fmt.Errorf("insufficient machines available with %v Cores, %vGB RAM, %vGB Disk", expected.NumCPU, expected.RAM, expected.Disk)

	for seen := 0; need > 0; seen++ {
		// not enough machines left to fill requirement
		if len(*resources)-seen < need {
			return "", errMsg
		}

		resource := (*resources)[seen]
		// check that the machine has the required resources
		if resource.NumCPU < expected.NumCPU || resource.RAM < expected.RAM || resource.Disk < expected.Disk {
			continue
		}

		// if pool is specified, check that the machine is in the required pool
		if expected.Pool != "" && resource.Pool != expected.Pool {
			continue
		}

		// if tags are specified, check that the machine has all the required tags
		if len(expected.Tags) > 0 {
			for _, tag := range expected.Tags {
				if !containsString(resource.Tags, tag) {
					continue
				}
			}
		}

		// once all checks pass, remove the machine from the list of available machines
		*resources = append((*resources)[:seen], (*resources)[seen+1:]...)
		need--
		// decrement seen to account for the removed resource
		seen--
	}

	return fmt.Sprintf("Found %d machine(s) available with %v Cores, %vGB RAM, %vGB Disk", expected.NumMachines, expected.NumCPU, expected.RAM, expected.Disk), nil
}

func formatMachines(machines []entity.Machine) []resource {
	resources := make([]resource, 0)
	for _, m := range machines {
		resources = append(resources, resource{
			NumCPU: m.CPUCount,
			RAM:    int(m.Memory) / 1024,
			Disk:   int(m.Storage) / 1000,
			Pool:   m.Pool.Name,
			Tags:   m.TagNames,
		})
	}

	// sort resources by ascending CPU count
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].NumCPU < resources[j].NumCPU
	})
	return resources
}
