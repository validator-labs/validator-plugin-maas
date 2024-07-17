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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MaasValidatorSpec defines the desired state of MaasValidator
type MaasValidatorSpec struct {
	Auth Auth   `json:"auth" yaml:"auth"`
	Host string `json:"host" yaml:"host"`

	OSImageRules              []OSImageRule              `json:"osImageRules,omitempty" yaml:"osImageRules,omitempty"`
	InternalDNSRules          []InternalDNSRule          `json:"internalDNSRules,omitempty" yaml:"internalDNSRules,omitempty"`
	ExternalDNSRules          []ExternalDNSRule          `json:"externalDNSRules,omitempty" yaml:"externalDNSRules,omitempty"`
	ResourceAvailabilityRules []ResourceAvailabilityRule `json:"resourceAvailabilityRules,omitempty" yaml:"resourceAvailabilityRules,omitempty"`
}

// ResultCount returns the number of validation results expected for an MaasValidatorSpec.
func (s MaasValidatorSpec) ResultCount() int {
	return len(s.OSImageRules) + len(s.InternalDNSRules) + len(s.ExternalDNSRules) + len(s.ResourceAvailabilityRules)
}

// Auth provides authentication information for the MaaS Instance
type Auth struct {
	// +kubebuilder:validation:Optional
	SecretName string `json:"secretName" yaml:"secretName"`
	// +kubebuilder:validation:Optional
	TokenKey string `json:"token" yaml:"token"`
}

// OSImageRule defines a rule for validating one or more OS images
type OSImageRule struct {
	// Unique name for the rule
	Name string `json:"name" yaml:"name"`
	// The list of OS images to validate
	OSImages []OSImage `json:"osImages" yaml:"osImages"`
}

// OSImage defines one OS image
type OSImage struct {
	// The name of the bootable image
	OSName string `json:"osName" yaml:"osName"`
	// OS Architecture
	Architecture string `json:"osArch" yaml:"osArch"`
}

// InternalDNSRule provides rules for the internal DNS server
type InternalDNSRule struct {
	// The domain name for the internal DNS server
	MaasDomain string `json:"maasDomain" yaml:"maasDomain"`
	// The DNS records for the internal DNS server
	DNSRecords []DNSRecord `json:"dnsRecords" yaml:"dnsRecords"`
}

// DNSRecord provides an internal DNS record
type DNSRecord struct {
	// The hostname for the DNS record
	Hostname string `json:"hostname" yaml:"hostname"`
	// The IP address for the DNS record
	IP string `json:"ip" yaml:"ip"`
	// The type of DNS record
	Type string `json:"type" yaml:"type"`
	// Optional Time To Live (TTL) for the DNS record
	TTL int `json:"ttl,omitempty" yaml:"ttl,omitempty"`
}

// ExternalDNSRule provides rules for validating the external DNS server
type ExternalDNSRule struct {
	// Unique name for the rule
	Name string `json:"name" yaml:"name"`
	// Whether the external DNS server is enabled
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// ResourceAvailabilityRule provides rules for validating resource availability
type ResourceAvailabilityRule struct {
	// Unique name for the rule
	Name string `json:"name" yaml:"name"`
	// The list of resources to validate
	Resources []Resource `json:"resources" yaml:"resources"`
}

// Resource defines a compute resource
type Resource struct {
	// Availability Zone
	AZ string `json:"az" yaml:"az"`
	// Minimum desired number of machines
	NumMachines int `json:"numMachines" yaml:"numMachines"`
	// Minimum CPU cores per machine
	NumCPU int `json:"numCPU" yaml:"numCPU"`
	// Minimum RAM per machine in GB
	NumRAM int `json:"numRAM" yaml:"numRAM"`
	// Minimum Disk space per machine in GB
	NumDisk int `json:"numDisk" yaml:"numDisk"`
	// Optional machine pool
	Pool string `json:"pool,omitempty" yaml:"pool,omitempty"`
	// Optional machine labels
	Labels []string `json:"labels,omitempty" yaml:"labels,omitempty"`
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
