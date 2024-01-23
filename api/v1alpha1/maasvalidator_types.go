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
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="MaasInstanceRules must have a unique Host",rule="self.all(e, size(self.filter(x, x.host == e.host)) == 1)"
	MaasInstanceRules []MaasInstanceRule `json:"ociRegistryRules,omitempty" yaml:"ociRegistryRules,omitempty"`
}

func (s MaasValidatorSpec) ResultCount() int {
	return len(s.MaasInstanceRules)
}

type MaasInstanceRule struct {
	// Host is a reference to the host URL of an OCI compliant registry
	Host string `json:"host" yaml:"host"`

	// Artifacts is a slice of artifacts in the host registry that should be validated.
	Artifacts []Artifact `json:"artifacts,omitempty" yaml:"artifacts,omitempty"`

	// Auth provides authentication information for the registry
	Auth Auth `json:"auth,omitempty" yaml:"auth,omitempty"`

	// CaCert is the base64 encoded CA Certificate
	CaCert string `json:"caCert,omitempty" yaml:"caCert,omitempty"`
}

func (r MaasInstanceRule) Name() string {
	return r.Host
}

type Artifact struct {
	// Ref is the path to the artifact in the host registry that should be validated.
	// An individual artifact can take any of the following forms:
	// <repository-path>/<artifact-name>
	// <repository-path>/<artifact-name>:<tag>
	// <repository-path>/<artifact-name>@<digest>
	//
	// When no tag or digest are specified, the default tag "latest" is used.
	Ref string `json:"ref" yaml:"ref"`

	// Download specifies whether a download attempt should be made for the artifact
	Download bool `json:"download,omitempty" yaml:"download,omitempty"`
}

type Auth struct {
	SecretName string `json:"secretName" yaml:"secretName"`
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
