// Package res contains the logic for validating resource rules
package res

import (
	"fmt"
	"slices"
	"sort"

	"github.com/canonical/gomaasclient/api"
	"github.com/canonical/gomaasclient/entity"
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/utils"
	"github.com/validator-labs/validator-plugin-maas/pkg/constants"
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
func (s *ResourceRulesService) ReconcileMaasInstanceResourceRule(rule v1alpha1.ResourceAvailabilityRule, seen map[string]bool) (*types.ValidationRuleResult, error) {
	errs := make([]error, 0)
	details := make([]string, 0)

	vr := utils.BuildValidationResult(rule.Name)

	// do not process an AZ more than once
	if seen[rule.AZ] {
		errs = append(errs, fmt.Errorf("availability zone %s referenced in a previous rule; AZ resource requirements must be defined in a single rule", rule.AZ))
	} else {
		errs, details = s.findMachineResources(rule)
	}

	utils.UpdateResult(vr, errs, constants.ErrImageNotFound, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// findMachineResources checks whether the rule is satisfied by the machines in the availability zone
func (s *ResourceRulesService) findMachineResources(rule v1alpha1.ResourceAvailabilityRule) ([]error, []string) {
	errs := make([]error, 0)
	details := make([]string, 0)

	resources, err := s.getAvailableMAASResources(rule)
	if err != nil {
		errs = append(errs, err)
		return errs, details
	}

	// each rule may have multiple resource checks for 1 AZ
	for _, rr := range rule.Resources {
		var matching []resource
		resources, matching, err = s.compareResources(resources, rr)
		if err != nil {
			errs = append(errs, err)
		}

		if len(matching) == rr.NumMachines {
			details = append(details, fmt.Sprintf("Found %d machine(s) available with %v Cores, %vGB RAM, %vGB Disk", len(matching), rr.NumCPU, rr.RAM, rr.Disk))
		}
	}

	return errs, details
}

func (s *ResourceRulesService) compareResources(resources []resource, expected v1alpha1.Resource) ([]resource, []resource, error) {
	if len(resources) < expected.NumMachines {
		return resources, nil, fmt.Errorf("not enough resources available in az: have: %v, need: %v", len(resources), expected.NumMachines)
	}

	matching := make([]resource, 0)
	remaining := make([]resource, 0)

	for _, r := range resources {
		ok := true

		if r.NumCPU < expected.NumCPU || r.RAM < expected.RAM || r.Disk < expected.Disk {
			ok = false
		}
		// if pool is specified, check that the machine is in the required pool
		if expected.Pool != "" && r.Pool != expected.Pool {
			ok = false
		}
		// if tags are specified, check that the machine has all the required tags
		if len(expected.Tags) > 0 {
			for _, tag := range expected.Tags {
				if !slices.Contains(r.Tags, tag) {
					ok = false
				}
			}
		}

		if ok {
			matching = append(matching, r)
		} else {
			remaining = append(remaining, r)
		}

		if len(matching) == expected.NumMachines {
			remaining = append(remaining, resources[len(matching):]...)
			return remaining, matching, nil
		}
	}

	err := fmt.Errorf("insufficient machines available with %v Cores, %vGB RAM, %vGB Disk. %v/%v available", expected.NumCPU, expected.RAM, expected.Disk, len(matching), expected.NumMachines)
	return resources, matching, err
}

func (s *ResourceRulesService) getAvailableMAASResources(rule v1alpha1.ResourceAvailabilityRule) ([]resource, error) {
	machines, err := s.api.Get(&entity.MachinesParams{Zone: []string{rule.AZ}, Status: []string{"ready"}})
	if err != nil {
		return nil, fmt.Errorf("error retrieving machines in availability zone %s: %w", rule.AZ, err)
	}

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
	return resources, nil
}
