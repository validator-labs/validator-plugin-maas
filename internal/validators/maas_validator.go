package validators

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	gomaasclient "github.com/maas/gomaasclient/client"
	"github.com/maas/gomaasclient/entity"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/spectrocloud-labs/validator-plugin-maas/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-maas/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util"
)

const errMsg string = "failed to validate rule"

type Resolver interface {
	LookupHost(ctx context.Context, nameserver string) (addrs []string, err error)
	// LookupHost(host string) (addrs []string, err error) {
}

type MaasRuleService struct {
	apiclient MaaSAPIClient
	resolvers []Resolver
	rr        map[v1alpha1.Nameserver]Resolver
}

type MaaSAPIClient interface {
	ListOSImages() ([]entity.BootResource, error)
	ListDNSServers() ([]v1alpha1.Nameserver, error)
}

type MaaSAPI struct {
	Client *gomaasclient.Client
}

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

// list the configured DNS servers in maas server configuration
// https://github.com/maas/gomaasclient/blob/master/api/maas_server.go
// Using the API: http://<your-server>/api/2.0/maas/op-get_config?name=upstream_dns
// will return a list of space separated Name servers: "10.10.128.8 8.8.4.4"
func (m *MaaSAPI) ListDNSServers() ([]v1alpha1.Nameserver, error) {
	if m.Client != nil {
		nameservers, err := m.Client.MAASServer.Get("upstream_dns")
		if err != nil {
			return make([]v1alpha1.Nameserver, 0), err
		}

		nameserversString := strings.Split(string(nameservers), " ")
		nameservsersMaas := make([]v1alpha1.Nameserver, len(nameserversString))
		for i, ns := range nameserversString {
			nameservsersMaas[i] = v1alpha1.Nameserver(ns)
		}
		return nameservsersMaas, nil
	}
	return make([]v1alpha1.Nameserver, 0), nil
}

func NewMaasRuleService(apiclient MaaSAPIClient) *MaasRuleService {

	mrs := &MaasRuleService{
		apiclient: apiclient,
		resolvers: make([]Resolver, 0),
		rr:        make(map[v1alpha1.Nameserver]Resolver),
	}

	return mrs
}

func (s *MaasRuleService) SetResolvers(resolvers map[v1alpha1.Nameserver]Resolver) {
	s.rr = resolvers
}

// ReconcileMaasInstanceRule reconciles a MaaS instance rule from the MaasValidator config
func (s *MaasRuleService) ReconcileMaasInstanceImageRules(rules v1alpha1.MaasInstanceRules) (*vapitypes.ValidationResult, error) {

	vr := buildValidationResult(rules)

	brs, err := s.ListOSImages()
	if err != nil {
		return vr, err
	}

	errs, details := findBootResources(rules.OSImages, brs)

	s.updateResult(vr, errs, errMsg, rules.Name, details...)

	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

func (s *MaasRuleService) ReconsileMaasInstanceExtDNSRules(rule v1alpha1.MaasExternalDNSRule) (*vapitypes.ValidationResult, error) {

	vr := buildValidationResult(rule)

	extDNS, err := s.apiclient.ListDNSServers()
	if err != nil {
		return vr, err
	}

	resolvers := make(map[v1alpha1.Nameserver]Resolver, 0)

	for _, ns := range extDNS {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(10000),
				}
				return d.DialContext(ctx, network, fmt.Sprintf("%s:53", ns))
			},
		}
		r.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:53", ns))
		}
		resolvers[ns] = r
	}
	s.SetResolvers(resolvers)
	errs, details := assertExternalDNS(s.rr)
	s.updateResult(vr, errs, errMsg, getType(rule), details...)
	if len(errs) > 0 {
		return vr, errs[0]
	}
	return vr, nil
}

func getType(rule interface{}) string {
	t := reflect.TypeOf(rule)
	return t.Name()
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule interface{}) *types.ValidationResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", constants.MaasInstance)
	latestCondition.ValidationRule = getType(rule)
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
	convertedSet := mapset.NewSet[v1alpha1.OSImage](converted...)
	imgRulesSet := mapset.NewSet[v1alpha1.OSImage](imgRules...)

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

func assertExternalDNS(resolvers map[v1alpha1.Nameserver]Resolver) (errs []error, details []string) {
	errs = make([]error, 0)
	details = make([]string, 0)

	if len(resolvers) == 0 {
		errs = append(errs, errors.New(errMsg))
		details = append(details, "No external nameservers found")
		return errs, details
	}

	for ns, rslvr := range resolvers {
		if ok, err := dnsLookUpIsWorking(ns, rslvr); !ok {
			errs = append(errs, err)
			details = append(details, fmt.Sprintf("Failed to resolve DNS with %s", ns))
		}
	}
	return errs, details
}

func dnsLookUpIsWorking(nameserver v1alpha1.Nameserver, resolver Resolver) (bool, error) {
	ip, err := resolver.LookupHost(context.Background(), "www.google.com")
	if err != nil {
		return false, err
	}
	if len(ip) > 1 {
		return true, nil
	}

	return false, fmt.Errorf("failed DNS resolution with %s", string(nameserver))
}
