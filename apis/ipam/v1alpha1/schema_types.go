/*
Copyright 2021 Nddr.

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
	nddov1 "github.com/yndd/nddo-runtime/apis/common/v1"
)

// NddrIpam struct
type NddrIpam struct {
	Ipam *NddrIpamIpam `json:"ipam,omitempty"`
}

// NddrIpamIpam struct
type NddrIpamIpam struct {
	Name       *string `json:"-name,omitempty"`
	TenantName *string `json:"tenant-name,omitempty"`
	VpcName    *string `json:"vpc-name,omitempty"`
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState      *string                        `json:"admin-state,omitempty"`
	Description     *string                        `json:"description,omitempty"`
	NetworkInstance []*NddrIpamIpamNetworkInstance `json:"network-instance,omitempty"`
	State           *NddrIpamIpamState             `json:"state,omitempty"`
}

// NddrIpamIpamNetworkInstance struct
type NddrIpamIpamNetworkInstance struct {
	AdminState         *string                                 `json:"admin-state,omitempty"`
	AllocationStrategy *string                                 `json:"allocation-strategy,omitempty"`
	Description        *string                                 `json:"description,omitempty"`
	IpAddress          []*NddrIpamIpamNetworkInstanceIpAddress `json:"ip-address,omitempty"`
	IpPrefix           []*NddrIpamIpamNetworkInstanceIpPrefix  `json:"ip-prefix,omitempty"`
	IpRange            []*NddrIpamIpamNetworkInstanceIpRange   `json:"ip-range,omitempty"`
	Name               *string                                 `json:"name,omitempty"`
	State              *NddrIpamIpamNetworkInstanceState       `json:"state,omitempty"`
	Tag                []*nddov1.Tag                           `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpAddress struct
type NddrIpamIpamNetworkInstanceIpAddress struct {
	Address     *string                                    `json:"address"`
	AdminState  *string                                    `json:"admin-state,omitempty"`
	Description *string                                    `json:"description,omitempty"`
	DnsName     *string                                    `json:"dns-name,omitempty"`
	NatInside   *string                                    `json:"nat-inside,omitempty"`
	NatOutside  *string                                    `json:"nat-outside,omitempty"`
	State       *NddrIpamIpamNetworkInstanceIpAddressState `json:"state,omitempty"`
	Tag         []*nddov1.Tag                              `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpAddressState struct
type NddrIpamIpamNetworkInstanceIpAddressState struct {
	//LastUpdate *string                                              `json:"last-update,omitempty"`
	//Origin     *string                                              `json:"origin,omitempty"`
	IpPrefix []*NddrIpamIpamNetworkInstanceIpAddressStateIpPrefix `json:"ip-prefix,omitempty"`
	IpRange  []*NddrIpamIpamNetworkInstanceIpAddressStateIpRange  `json:"ip-range,omitempty"`
	Reason   *string                                              `json:"reason,omitempty"`
	Status   *string                                              `json:"status,omitempty"`
	Tag      []*nddov1.Tag                                        `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpAddressStateIpPrefix struct
type NddrIpamIpamNetworkInstanceIpAddressStateIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddrIpamIpamNetworkInstanceIpAddressStateIpRange struct
type NddrIpamIpamNetworkInstanceIpAddressStateIpRange struct {
	End   *string `json:"end"`
	Start *string `json:"start"`
}

// NddrIpamIpamNetworkInstanceIpPrefix struct
type NddrIpamIpamNetworkInstanceIpPrefix struct {
	AdminState  *string `json:"admin-state,omitempty"`
	Description *string `json:"description,omitempty"`
	Pool        *bool   `json:"pool,omitempty"`
	Prefix      *string `json:"prefix"`
	//RirName     *string                                   `json:"rir-name,omitempty"`
	State *NddrIpamIpamNetworkInstanceIpPrefixState `json:"state,omitempty"`
	Tag   []*nddov1.Tag                             `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpPrefixState struct
type NddrIpamIpamNetworkInstanceIpPrefixState struct {
	Adresses *uint32                                        `json:"adresses,omitempty"`
	Child    *NddrIpamIpamNetworkInstanceIpPrefixStateChild `json:"child,omitempty"`
	//LastUpdate *string                                         `json:"last-update,omitempty"`
	//Origin     *string                                         `json:"origin,omitempty"`
	Parent *NddrIpamIpamNetworkInstanceIpPrefixStateParent `json:"parent,omitempty"`
	Reason *string                                         `json:"reason,omitempty"`
	Status *string                                         `json:"status,omitempty"`
	Tag    []*nddov1.Tag                                   `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpPrefixStateChild struct
type NddrIpamIpamNetworkInstanceIpPrefixStateChild struct {
	IpPrefix []*NddrIpamIpamNetworkInstanceIpPrefixStateChildIpPrefix `json:"ip-prefix,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpPrefixStateChildIpPrefix struct
type NddrIpamIpamNetworkInstanceIpPrefixStateChildIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddrIpamIpamNetworkInstanceIpPrefixStateParent struct
type NddrIpamIpamNetworkInstanceIpPrefixStateParent struct {
	IpPrefix []*NddrIpamIpamNetworkInstanceIpPrefixStateParentIpPrefix `json:"ip-prefix,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpPrefixStateParentIpPrefix struct
type NddrIpamIpamNetworkInstanceIpPrefixStateParentIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddrIpamIpamNetworkInstanceIpRange struct
type NddrIpamIpamNetworkInstanceIpRange struct {
	AdminState  *string                                  `json:"admin-state,omitempty"`
	Description *string                                  `json:"description,omitempty"`
	End         *string                                  `json:"end"`
	Start       *string                                  `json:"start"`
	State       *NddrIpamIpamNetworkInstanceIpRangeState `json:"state,omitempty"`
	Tag         []*nddov1.Tag                            `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRangeState struct
type NddrIpamIpamNetworkInstanceIpRangeState struct {
	//LastUpdate *string                                        `json:"last-update,omitempty"`
	//Origin     *string                                        `json:"origin,omitempty"`
	Parent *NddrIpamIpamNetworkInstanceIpRangeStateParent `json:"parent,omitempty"`
	Reason *string                                        `json:"reason,omitempty"`
	Size   *uint32                                        `json:"size,omitempty"`
	Status *string                                        `json:"status,omitempty"`
	Tag    []*nddov1.Tag                                  `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRangeStateParent struct
type NddrIpamIpamNetworkInstanceIpRangeStateParent struct {
	IpPrefix []*NddrIpamIpamNetworkInstanceIpRangeStateParentIpPrefix `json:"ip-prefix,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRangeStateParentIpPrefix struct
type NddrIpamIpamNetworkInstanceIpRangeStateParentIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddrIpamIpamNetworkInstanceState struct
type NddrIpamIpamNetworkInstanceState struct {
	//LastUpdate *string                                `json:"last-update,omitempty"`
	//Origin     *string                                `json:"origin,omitempty"`
	Reason *string       `json:"reason,omitempty"`
	Status *string       `json:"status,omitempty"`
	Tag    []*nddov1.Tag `json:"tag,omitempty"`
}

// NddrIpamIpamState struct
type NddrIpamIpamState struct {
	//LastUpdate *string                 `json:"last-update,omitempty"`
	//Origin     *string                 `json:"origin,omitempty"`
	Reason *string       `json:"reason,omitempty"`
	Status *string       `json:"status,omitempty"`
	Tag    []*nddov1.Tag `json:"tag,omitempty"`
}

// Root is the root of the schema
type Root struct {
	IpamNddrIpam *NddrIpam `json:"Nddr-ipam,omitempty"`
}
