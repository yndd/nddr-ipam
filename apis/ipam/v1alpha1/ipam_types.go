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
	// IpamFinalizer is the name of the finalizer added to
	// Ipam to block delete operations until the physical node can be
	// deprovisioned.
	IpamFinalizer string = "ipam.ipam.nddr.yndd.io"
)

// Ipam struct
type IpamIpam struct {
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState  *string `json:"admin-state,omitempty"`
	Description *string `json:"description,omitempty"`
	//Rir []*IpamRir `json:"rir,omitempty"`
}

// A IpamSpec defines the desired state of a Ipam.
type IpamSpec struct {
	//nddv1.ResourceSpec `json:",inline"`
	Ipam *IpamIpam `json:"ipam,omitempty"`
}

// A IpamStatus represents the observed state of a Ipam.
type IpamStatus struct {
	nddv1.ConditionedStatus `json:",inline"`
	OrganizationName        *string       `json:"organization-name,omitempty"`
	IpamName                *string       `json:"ipam-name,omitempty"`
	Ipam                    *NddrIpamIpam `json:"ipam,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpam is the Schema for the Ipam API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.organization-name"
// +kubebuilder:printcolumn:name="iPAM",type="string",JSONPath=".status.ipam-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type Ipam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamSpec   `json:"spec,omitempty"`
	Status IpamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamList contains a list of Ipams
type IpamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ipam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ipam{}, &IpamList{})
}

// Ipam type metadata.
var (
	IpamKindKind         = reflect.TypeOf(Ipam{}).Name()
	IpamGroupKind        = schema.GroupKind{Group: Group, Kind: IpamKindKind}.String()
	IpamKindAPIVersion   = IpamKindKind + "." + GroupVersion.String()
	IpamGroupVersionKind = GroupVersion.WithKind(IpamKindKind)
)
