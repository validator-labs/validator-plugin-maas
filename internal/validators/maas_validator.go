package validators

import (
	"errors"
	"fmt"

	gomaasclient "github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	mapset "github.com/deckarep/golang-set/v2"

	"github.com/validator-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
)

const errMsg string = "failed to validate rule"

// MaasRuleService is a service for reconciling Maas rules via its MaaSAPIClient
type MaasRuleService struct {
	apiclient MaaSAPIClient
}

// MaasAPIClient is an interface for interacting with the Maas API
type MaaSAPIClient interface {
	ListOSImages() ([]entity.BootResource, error)
	ListDNSServers() ([]entity.DNSResource, error)
}

// MaaSAPI is a struct for which containts the Maas client
type MaaSAPI struct {
	Client *gomaasclient.Client
}

// ListOSImages returns a list of OS images from the Maas API
func (m *MaaSAPI) ListOSImages() ([]entity.BootResource, error) {
	if m.Client != nil {
		images, err := m.Client.BootResources.Get(&entity.BootResourcesReadParams{})
		if err != nil {
			return make([]entity.BootResource, 0), err
		}
		return images, nil
	}
	return make([]entity.BootResource, 0), nil
}

// ListDNSServers returns a list of DNS servers from the Maas API
func (m *MaaSAPI) ListDNSServers() ([]entity.DNSResource, error) {
	if m.Client != nil {
		dnsresources, err := m.Client.DNSResources.Get()
		if err != nil {
			return make([]entity.DNSResource, 0), err
		}
		return dnsresources, nil
	}
	return make([]entity.DNSResource, 0), nil
}

// NewMaasRuleService returns a MaasRuleService
func NewMaasRuleService(apiclient MaaSAPIClient) *MaasRuleService {
	return &MaasRuleService{
		apiclient: apiclient,
	}
}

// ReconcileMaasInstanceRule reconciles a MaaS instance rule from the MaasValidator config
func (s *MaasRuleService) ReconcileMaasInstanceImageRules(rules v1alpha1.MaasInstanceRules) (*types.ValidationRuleResult, error) {

	vr := buildValidationResult(rules)

	brs, err := s.ListOSImages()
	if err != nil {
		return vr, err
	}

	errs, details := findBootResources(rules.OSImages, brs)

	s.updateResult(vr, errs, errMsg, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rules v1alpha1.MaasInstanceRules) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.MaasInstance)
	latestCondition.ValidationRule = rules.Name
	latestCondition.ValidationType = constants.MaasInstance
	return &types.ValidationRuleResult{Condition: &latestCondition, State: &state}
}

func (s *MaasRuleService) updateResult(vr *types.ValidationRuleResult, errs []error, errMsg string, details ...string) {
	if len(errs) > 0 {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Message = errMsg
		for _, err := range errs {
			vr.Condition.Failures = append(vr.Condition.Failures, err.Error())
		}
	}
	vr.Condition.Details = append(vr.Condition.Details, details...)
}

// ListOSImages returns a list of OS images from the Maas API
func (s *MaasRuleService) ListOSImages() ([]entity.BootResource, error) {
	images, err := s.apiclient.ListOSImages()
	if err != nil {
		return nil, err
	}
	return images, nil
}

func convertBootResourceToOSImage(images []entity.BootResource) []v1alpha1.OSImage {
	converted := make([]v1alpha1.OSImage, len(images))
	for i, img := range images {
		converted[i] = v1alpha1.OSImage{
			Name:         img.Name,
			Architecture: img.Architecture,
		}
	}
	return converted
}

func findBootResources(imgRules []v1alpha1.OSImage, images []entity.BootResource) (errs []error, details []string) {
	errs = make([]error, 0)
	details = make([]string, 0)

	converted := convertBootResourceToOSImage(images)
	convertedSet := mapset.NewSet(converted...)
	imgRulesSet := mapset.NewSet(imgRules...)

	if imgRulesSet.IsSubset(convertedSet) {
		return errs, details
	}

	diffSet := imgRulesSet.Difference(convertedSet)

	diffSetIt := diffSet.Iterator()

	for img := range diffSetIt.C {
		errs = append(errs, errors.New(errMsg))
		details = append(details, fmt.Sprintf("OS image %s with arch %s was not found", img.Name, img.Architecture))
	}

	return errs, details
}
