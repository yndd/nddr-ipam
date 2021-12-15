/*
Copyright 2021 NDDO.

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

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// AllocFinalizer is the name of the finalizer added to
	// Alloc to block delete operations until all resources have been finalized
	AllocFinalizer string = "alloc.ipam.nddr.yndd.io"
)

// NddrIpamAlloc struct
type NddrIpamAlloc struct {
	IpamAlloc `json:",inline"`
	State     *NddrAllocState `json:"state,omitempty"`
}

// NddrAllocState struct
type NddrAllocState struct {
	IpPrefix *string `json:"ip-prefix,omitempty"`
	//ExpiryTime *string `json:"expiry-time,omitempty"`
}

// IpamAlloc struct
type IpamAlloc struct {
	// +kubebuilder:validation:Enum=`ipv4`;`ipv6`
	// +kubebuilder:default:="ipv4"
	AddressFamily *string `json:"address-family,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=128
	PrefixLength *uint32                  `json:"prefix-length,omitempty"`
	Selector     []*IpamAllocSelectorTag  `json:"selector,omitempty"`
	SourceTag    []*IpamAllocSourceTagTag `json:"source-tag,omitempty"`
}

type IpamAllocSelectorTag struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type IpamAllocSourceTagTag struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

// A AllocSpec defines the desired state of a Alloc.
type AllocSpec struct {
	//nddv1.ResourceSpec `json:",inline"`
	IpamName            *string    `json:"ipam-name,omitempty"`
	NetworkInstanceName *string    `json:"network-instance-name,omitempty"`
	Alloc               *IpamAlloc `json:"alloc,omitempty"`
}

// A AllocStatus represents the observed state of a Alloc.
type AllocStatus struct {
	nddv1.ConditionedStatus `json:",inline"`
	Alloc                   *NddrIpamAlloc `json:"alloc,omitempty"`
}

// +kubebuilder:object:root=true

// Alloc is the Schema for the Alloc API
// +kubebuilder:subresource:status
type Alloc struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AllocSpec   `json:"spec,omitempty"`
	Status AllocStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AllocList contains a list of Allocs
type AllocList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Alloc `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Alloc{}, &AllocList{})
}

// Alloc type metadata.
var (
	AllocKindKind         = reflect.TypeOf(Alloc{}).Name()
	AllocGroupKind        = schema.GroupKind{Group: Group, Kind: AllocKindKind}.String()
	AllocKindAPIVersion   = AllocKindKind + "." + GroupVersion.String()
	AllocGroupVersionKind = GroupVersion.WithKind(AllocKindKind)
)
