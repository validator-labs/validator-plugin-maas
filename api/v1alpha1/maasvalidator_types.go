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
	MaasInstanceRules `json:"MaasInstanceRules,omitempty" yaml:"MaasInstanceRules,omitempty"`
	MaasInstance      `json:"MaasInstance" yaml:"MaasInstance"`
}

func (s MaasValidatorSpec) ResultCount() int {
	return len(s.MaasInstanceRules.OSImages)
}

// MaasInstance describes the MaaS host
type MaasInstance struct {
	// Host is the URL for your MaaS instance
	Host string `json:"host" yaml:"host"`
	Auth Auth   `json:"auth" yaml:"auth"`
}

type MaasInstanceRules struct {
	// Unique rule name
	Name string `json:"name" yaml:"name"`
	// OSImages is a list of bootable os images
	OSImages []OSImage `json:"bootable-images,omitempty" yaml:"bootable-images,omitempty"`

	// Auth provides authentication information for the MaaS Instance
	Auth Auth `json:"auth" yaml:"auth"`
}

type OSImage struct {
	// The name of the bootable image
	Name string `json:"name" yaml:"name"`
	// OS Architecture
	Architecture string `json:"os-arch" yaml:"os-arch"`
}

type Auth struct {
	// +kubebuilder:validation:Optional
	SecretName string `json:"secretName" yaml:"secretName"`
	// +kubebuilder:validation:Optional
	TokenKey string `json:"token" yaml:"token"`
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
