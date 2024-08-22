//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Auth) DeepCopyInto(out *Auth) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Auth.
func (in *Auth) DeepCopy() *Auth {
	if in == nil {
		return nil
	}
	out := new(Auth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecord) DeepCopyInto(out *DNSRecord) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecord.
func (in *DNSRecord) DeepCopy() *DNSRecord {
	if in == nil {
		return nil
	}
	out := new(DNSRecord)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSResource) DeepCopyInto(out *DNSResource) {
	*out = *in
	if in.DNSRecords != nil {
		in, out := &in.DNSRecords, &out.DNSRecords
		*out = make([]DNSRecord, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSResource.
func (in *DNSResource) DeepCopy() *DNSResource {
	if in == nil {
		return nil
	}
	out := new(DNSResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Image.
func (in *Image) DeepCopy() *Image {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageRule) DeepCopyInto(out *ImageRule) {
	*out = *in
	out.ManuallyNamed = in.ManuallyNamed
	if in.Images != nil {
		in, out := &in.Images, &out.Images
		*out = make([]Image, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageRule.
func (in *ImageRule) DeepCopy() *ImageRule {
	if in == nil {
		return nil
	}
	out := new(ImageRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalDNSRule) DeepCopyInto(out *InternalDNSRule) {
	*out = *in
	out.AutomaticallyNamed = in.AutomaticallyNamed
	if in.DNSResources != nil {
		in, out := &in.DNSResources, &out.DNSResources
		*out = make([]DNSResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalDNSRule.
func (in *InternalDNSRule) DeepCopy() *InternalDNSRule {
	if in == nil {
		return nil
	}
	out := new(InternalDNSRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MaasValidator) DeepCopyInto(out *MaasValidator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MaasValidator.
func (in *MaasValidator) DeepCopy() *MaasValidator {
	if in == nil {
		return nil
	}
	out := new(MaasValidator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MaasValidator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MaasValidatorList) DeepCopyInto(out *MaasValidatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MaasValidator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MaasValidatorList.
func (in *MaasValidatorList) DeepCopy() *MaasValidatorList {
	if in == nil {
		return nil
	}
	out := new(MaasValidatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MaasValidatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MaasValidatorSpec) DeepCopyInto(out *MaasValidatorSpec) {
	*out = *in
	out.Auth = in.Auth
	if in.ImageRules != nil {
		in, out := &in.ImageRules, &out.ImageRules
		*out = make([]ImageRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InternalDNSRules != nil {
		in, out := &in.InternalDNSRules, &out.InternalDNSRules
		*out = make([]InternalDNSRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.UpstreamDNSRules != nil {
		in, out := &in.UpstreamDNSRules, &out.UpstreamDNSRules
		*out = make([]UpstreamDNSRule, len(*in))
		copy(*out, *in)
	}
	if in.ResourceAvailabilityRules != nil {
		in, out := &in.ResourceAvailabilityRules, &out.ResourceAvailabilityRules
		*out = make([]ResourceAvailabilityRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MaasValidatorSpec.
func (in *MaasValidatorSpec) DeepCopy() *MaasValidatorSpec {
	if in == nil {
		return nil
	}
	out := new(MaasValidatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MaasValidatorStatus) DeepCopyInto(out *MaasValidatorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MaasValidatorStatus.
func (in *MaasValidatorStatus) DeepCopy() *MaasValidatorStatus {
	if in == nil {
		return nil
	}
	out := new(MaasValidatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceAvailabilityRule) DeepCopyInto(out *ResourceAvailabilityRule) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]Resource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceAvailabilityRule.
func (in *ResourceAvailabilityRule) DeepCopy() *ResourceAvailabilityRule {
	if in == nil {
		return nil
	}
	out := new(ResourceAvailabilityRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UpstreamDNSRule) DeepCopyInto(out *UpstreamDNSRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpstreamDNSRule.
func (in *UpstreamDNSRule) DeepCopy() *UpstreamDNSRule {
	if in == nil {
		return nil
	}
	out := new(UpstreamDNSRule)
	in.DeepCopyInto(out)
	return out
}
