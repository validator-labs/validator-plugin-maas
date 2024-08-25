/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/validator-labs/validator/pkg/plugins"
	"github.com/validator-labs/validator/pkg/validationrule"

	"github.com/validator-labs/validator-plugin-maas/pkg/constants"
)

// MaasValidatorSpec defines the desired state of MaasValidator
type MaasValidatorSpec struct {
	Auth Auth `json:"auth" yaml:"auth"`
	// Host is the URL for your MAAS Instance
	Host string `json:"host" yaml:"host"`

	ImageRules                []ImageRule                `json:"imageRules,omitempty" yaml:"imageRules,omitempty"`
	InternalDNSRules          []InternalDNSRule          `json:"internalDNSRules,omitempty" yaml:"internalDNSRules,omitempty"`
	UpstreamDNSRules          []UpstreamDNSRule          `json:"upstreamDNSRules,omitempty" yaml:"upstreamDNSRules,omitempty"`
	ResourceAvailabilityRules []ResourceAvailabilityRule `json:"resourceAvailabilityRules,omitempty" yaml:"resourceAvailabilityRules,omitempty"`
}

var _ plugins.PluginSpec = (*MaasValidatorSpec)(nil)

// PluginCode returns the MAAS validator's plugin code.
func (s MaasValidatorSpec) PluginCode() string {
	return constants.PluginCode
}

// ResultCount returns the number of validation results expected for an MaasValidatorSpec.
func (s MaasValidatorSpec) ResultCount() int {
	return len(s.ImageRules) + len(s.InternalDNSRules) + len(s.UpstreamDNSRules) + len(s.ResourceAvailabilityRules)
}

// Auth provides authentication information for the MAAS Instance
// Requires one of SecretName and TokenKey (if reading secret from Kubernetes) or APIToken (if reading secret from file)
type Auth struct {
	// SecretName is the name of the Kubernetes secret containing the MAAS API token
	SecretName string `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	// TokenKey is the key in the Kubernetes secret containing the MAAS API token
	TokenKey string `json:"tokenKey,omitempty" yaml:"tokenKey,omitempty"`
	// APIToken is the value of the MAAS API token
	APIToken string `json:"apiToken,omitempty" yaml:"apiToken,omitempty"`
}

// ImageRule defines a rule for validating one or more OS images
type ImageRule struct {
	validationrule.ManuallyNamed `json:",inline" yaml:",omitempty"`

	// Unique name for the rule
	RuleName string `json:"name" yaml:"name"`
	// The list of OS images to validate
	Images []Image `json:"images" yaml:"images"`
}

var _ validationrule.Interface = (*ImageRule)(nil)

// Name returns the name of the ImageRule.
func (r ImageRule) Name() string {
	return r.RuleName
}

// SetName sets the name of the ImageRule.
func (r *ImageRule) SetName(name string) {
	r.RuleName = name
}

// Image defines one OS image
type Image struct {
	// The name of the bootable image
	Name string `json:"name" yaml:"name"`
	// OS Architecture
	Architecture string `json:"architecture" yaml:"architecture"`
}

// InternalDNSRule provides rules for the internal DNS server
type InternalDNSRule struct {
	validationrule.AutomaticallyNamed `json:",inline" yaml:",omitempty"`

	// The domain name for the internal DNS server
	MaasDomain string `json:"maasDomain" yaml:"maasDomain"`
	// The DNS resources for the internal DNS server
	DNSResources []DNSResource `json:"dnsResources" yaml:"dnsResources"`
}

var _ validationrule.Interface = (*InternalDNSRule)(nil)

// Name returns the name of the InternalDNSRule.
func (r InternalDNSRule) Name() string {
	return r.MaasDomain
}

// DNSResource provides an internal DNS resource
type DNSResource struct {
	// The hostname for the DNS record
	FQDN string `json:"fqdn" yaml:"fqdn"`
	// The expected records for the FQDN
	DNSRecords []DNSRecord `json:"dnsRecords" yaml:"dnsRecords"`
}

// DNSRecord provides one DNS Resource Record
type DNSRecord struct {
	// The IP address for the DNS record
	IP string `json:"ip" yaml:"ip"`
	// The type of DNS record
	Type string `json:"type" yaml:"type"`
	// Optional Time To Live (TTL) for the DNS record
	TTL int `json:"ttl,omitempty" yaml:"ttl,omitempty"`
}

// UpstreamDNSRule provides rules for validating the external DNS server
type UpstreamDNSRule struct {
	validationrule.ManuallyNamed `json:",inline" yaml:",omitempty"`

	// Unique name for the rule
	RuleName string `json:"name" yaml:"name"`
	// The minimum expected number of upstream DNS servers
	NumDNSServers int `json:"numDNSServers" yaml:"numDNSServers"`
}

var _ validationrule.Interface = (*UpstreamDNSRule)(nil)

// Name returns the name of the UpstreamDNSRule.
func (r UpstreamDNSRule) Name() string {
	return r.RuleName
}

// SetName sets the name of the UpstreamDNSRule.
func (r *UpstreamDNSRule) SetName(name string) {
	r.RuleName = name
}

// ResourceAvailabilityRule provides rules for validating resource availability
type ResourceAvailabilityRule struct {
	validationrule.ManuallyNamed `json:",inline" yaml:",omitempty"`

	// Unique name for the rule
	RuleName string `json:"name" yaml:"name"`
	// The availability zone to validate
	AZ string `json:"az" yaml:"az"`
	// The list of resources to validate
	Resources []Resource `json:"resources" yaml:"resources"`
}

var _ validationrule.Interface = (*ResourceAvailabilityRule)(nil)

// Name returns the name of the RBAC rule.
func (r ResourceAvailabilityRule) Name() string {
	return r.RuleName
}

// SetName sets the name of the RBAC rule.
func (r *ResourceAvailabilityRule) SetName(name string) {
	r.RuleName = name
}

// Resource defines a compute resource
type Resource struct {
	// Minimum desired number of machines
	NumMachines int `json:"numMachines" yaml:"numMachines"`
	// Minimum CPU cores per machine
	NumCPU int `json:"numCPU" yaml:"numCPU"`
	// Minimum RAM per machine in GB
	RAM int `json:"ram" yaml:"ram"`
	// Minimum Disk space per machine in GB
	Disk int `json:"disk" yaml:"disk"`
	// Optional machine pool
	Pool string `json:"pool,omitempty" yaml:"pool,omitempty"`
	// Optional machine tags
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// MaasValidatorStatus defines the observed state of MaasValidator
type MaasValidatorStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MaasValidator is the Schema for the maasvalidators API
type MaasValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaasValidatorSpec   `json:"spec,omitempty"`
	Status MaasValidatorStatus `json:"status,omitempty"`
}

// GetKind returns the MAAS validator's kind.
func (v MaasValidator) GetKind() string {
	return reflect.TypeOf(v).Name()
}

// PluginCode returns the MAAS validator's plugin code.
func (v MaasValidator) PluginCode() string {
	return v.Spec.PluginCode()
}

// ResultCount returns the number of validation results expected for an MaasValidatorSpec.
func (v MaasValidator) ResultCount() int {
	return v.Spec.ResultCount()
}

//+kubebuilder:object:root=true

// MaasValidatorList contains a list of MaasValidator
type MaasValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaasValidator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaasValidator{}, &MaasValidatorList{})
}
