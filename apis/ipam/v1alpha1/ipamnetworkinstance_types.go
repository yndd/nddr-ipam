/*
Copyright 2021 nddr.

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
	// IpamTenantNetworkInstanceFinalizer is the name of the finalizer added to
	// IpamTenantNetworkInstance to block delete operations until the physical node can be
	// deprovisioned.
	IpamNetworkInstanceFinalizer string = "networkInstance.ipam.nddr.yndd.io"
)

// IpamTenantNetworkInstance struct
type IpamIpamNetworkInstance struct {
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState *string `json:"admin-state,omitempty"`
	// +kubebuilder:validation:Enum=`first-available`;`deterministic`
	// +kubebuilder:default:="first-available"
	AllocationStrategy *string `json:"allocation-strategy,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Description *string `json:"description,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	// +kubebuilder:default:="default"
	Name *string                       `json:"name,omitempty"`
	Tag  []*IpamIpamNetworkInstanceTag `json:"tag,omitempty"`
}

// IpamTenantNetworkInstanceTag struct
type IpamIpamNetworkInstanceTag struct {
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Key *string `json:"key"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Value *string `json:"value,omitempty"`
}

// A IpamSpec defines the desired state of a Ipam.
type IpamNetworkInstanceSpec struct {
	//nddv1.ResourceSpec `json:",inline"`
	IpamName            *string                  `json:"ipam-name,omitempty"`
	IpamNetworkInstance *IpamIpamNetworkInstance `json:"network-instance,omitempty"`
}

// A IpamStatus represents the observed state of a Ipam.
type IpamNetworkInstanceStatus struct {
	nddv1.ConditionedStatus `json:",inline"`
	IpamNetworkInstance     *NddrIpamIpamNetworkInstance `json:"network-instance,omitempty"`
}

// +kubebuilder:object:root=true

// IpamNetworkInstance is the Schema for the IpamNetworkInstance API
// +kubebuilder:subresource:status
type IpamNetworkInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamNetworkInstanceSpec   `json:"spec,omitempty"`
	Status IpamNetworkInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamNetworkInstanceList contains a list of IpamNetworkInstances
type IpamNetworkInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamNetworkInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamNetworkInstance{}, &IpamNetworkInstanceList{})
}

// IpamNetworkInstance type metadata.
var (
	IpamNetworkInstanceKindKind         = reflect.TypeOf(IpamNetworkInstance{}).Name()
	IpamNetworkInstanceGroupKind        = schema.GroupKind{Group: Group, Kind: IpamNetworkInstanceKindKind}.String()
	IpamNetworkInstanceKindAPIVersion   = IpamNetworkInstanceKindKind + "." + GroupVersion.String()
	IpamNetworkInstanceGroupVersionKind = GroupVersion.WithKind(IpamNetworkInstanceKindKind)
)
