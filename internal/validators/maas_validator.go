package validators

import (
	"errors"
	"fmt"

	gomaasclient "github.com/maas/gomaasclient/client"
	"github.com/maas/gomaasclient/entity"

	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util"
)

const errMsg string = "failed to validate rule"

type MaasRuleService struct {
	apiclient MaaSAPIClient
}

type MaaSAPIClient interface {
	ListOSImages() ([]entity.BootResource, error)
	ListDNSServers() ([]entity.DNSResource, error)
}

type MaaSAPI struct {
	Client *gomaasclient.Client
}

func (m *MaaSAPI) ListOSImages() ([]entity.BootResource, error) {
	images, _ := m.Client.BootResources.Get(&entity.BootResourcesReadParams{})
	return images, nil
}

func (m *MaaSAPI) ListDNSServers() ([]entity.DNSResource, error) {
	dnsresources, _ := m.Client.DNSResources.Get()
	return dnsresources, nil
}

func NewMaasRuleService(apiclient MaaSAPIClient) *MaasRuleService {
	return &MaasRuleService{
		apiclient: apiclient,
	}
}

// ReconcileMaasInstanceRule reconciles a MaaS instance rule from the MaasValidator config
func (s *MaasRuleService) ReconcileMaasInstanceRule(imgRule v1alpha1.OSImage) (*vapitypes.ValidationResult, error) {
	vr := buildValidationResult(imgRule)

	brs, err := s.ListOSImages()
	if err != nil {
		return vr, err
	}

	errs, details := findBootResources(imgRule, brs)

	s.updateResult(vr, errs, errMsg, imgRule.Name, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule v1alpha1.OSImage) *types.ValidationResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.MaasInstance)
	latestCondition.ValidationRule = rule.Name
	latestCondition.ValidationType = constants.MaasInstance
	return &types.ValidationResult{Condition: &latestCondition, State: &state}
}

func (s *MaasRuleService) updateResult(vr *types.ValidationResult, errs []error, errMsg, ruleName string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	for _, detail := range details {
		vr.Condition.Details = append(vr.Condition.Details, detail)
	}
}

func (s *MaasRuleService) ListOSImages() ([]entity.BootResource, error) {
	images, err := s.apiclient.ListOSImages()
	if err != nil {
		return nil, err
	}
	return images, nil
}

func findBootResources(imgRule v1alpha1.OSImage, images []entity.BootResource) (errs []error, details []string) {
	errs = make([]error, 0)
	details = make([]string, 0)

	var found bool = false
	for _, image := range images {
		if (imgRule.Name == image.Name) && (imgRule.Architecture == image.Architecture) {
			found = true
			break
		}
	}

	if !found {
		errs = append(errs, errors.New(errMsg))
		details = append(details, fmt.Sprintf("OS image %s with arch %s was not found", imgRule.Name, imgRule.Architecture))
	}

	return errs, details
}
