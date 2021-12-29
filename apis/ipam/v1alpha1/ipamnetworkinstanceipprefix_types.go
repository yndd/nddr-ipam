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
	nddov1 "github.com/yndd/nddo-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// IpamTenantNetworkInstanceIpPrefixFinalizer is the name of the finalizer added to
	// IpamTenantNetworkInstanceIpPrefix to block delete operations until the physical node can be
	// deprovisioned.
	IpamNetworkInstanceIpPrefixFinalizer string = "ipPrefix.ipam.nddr.yndd.io"
)

// IpamTenantNetworkInstanceIpPrefix struct
type IpamIpamNetworkInstanceIpPrefix struct {
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState *string `json:"admin-state,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Description *string `json:"description,omitempty"`
	Pool        *bool   `json:"pool,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Prefix *string `json:"prefix"`
	//RirName *string                                 `json:"rir-name,omitempty"`
	Tag []*nddov1.Tag `json:"tag,omitempty"`
}

// A IpamNetworkInstanceIpPrefixSpec defines the desired state of a IpamNetworkInstanceIpPrefix.
type IpamNetworkInstanceIpPrefixSpec struct {
	//nddv1.ResourceSpec `json:",inline"`
	IpamNetworkInstanceIpPrefix *IpamIpamNetworkInstanceIpPrefix `json:"ip-prefix,omitempty"`
}

// A IpamNetworkInstanceIpPrefixStatus represents the observed state of a IpamNetworkInstanceIpPrefix.
type IpamNetworkInstanceIpPrefixStatus struct {
	nddv1.ConditionedStatus     `json:",inline"`
	IpamName                    *string                              `json:"ipam-name,omitempty"`
	OrganizationName            *string                              `json:"organization-name,omitempty"`
	NetworkInstanceName         *string                              `json:"network-instance-name,omitempty"`
	IpPrefixName                *string                              `json:"ip-prefix-name,omitempty"`
	IpamNetworkInstanceIpPrefix *NddrIpamIpamNetworkInstanceIpPrefix `json:"ip-prefix,omitempty"`
}

// +kubebuilder:object:root=true

// IpamNetworkInstanceIpPrefix is the Schema for the IpamNetworkInstanceIpPrefix API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.organization-name"
// +kubebuilder:printcolumn:name="IPAM",type="string",JSONPath=".status.ipam-name"
// +kubebuilder:printcolumn:name="NI",type="string",JSONPath=".status.network-instance-name"
// +kubebuilder:printcolumn:name="PREFIX",type="string",JSONPath=".spec.ip-prefix.prefix"
// +kubebuilder:printcolumn:name="AF",type="string",JSONPath=".status.ip-prefix.state.tag[?(@.key=='address-family')].value"
// +kubebuilder:printcolumn:name="PURPOSE",type="string",JSONPath=".spec.ip-prefix.tag[?(@.key=='purpose')].value"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.ip-prefix.state.status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type IpamNetworkInstanceIpPrefix struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamNetworkInstanceIpPrefixSpec   `json:"spec,omitempty"`
	Status IpamNetworkInstanceIpPrefixStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamNetworkInstanceIpPrefixList contains a list of IpamNetworkInstanceIpPrefixes
type IpamNetworkInstanceIpPrefixList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamNetworkInstanceIpPrefix `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamNetworkInstanceIpPrefix{}, &IpamNetworkInstanceIpPrefixList{})
}

// IpamNetworkInstanceIpPrefix type metadata.
var (
	IpamNetworkInstanceIpPrefixKindKind         = reflect.TypeOf(IpamNetworkInstanceIpPrefix{}).Name()
	IpamNetworkInstanceIpPrefixGroupKind        = schema.GroupKind{Group: Group, Kind: IpamNetworkInstanceIpPrefixKindKind}.String()
	IpamNetworkInstanceIpPrefixKindAPIVersion   = IpamNetworkInstanceIpPrefixKindKind + "." + GroupVersion.String()
	IpamNetworkInstanceIpPrefixGroupVersionKind = GroupVersion.WithKind(IpamNetworkInstanceIpPrefixKindKind)
)
