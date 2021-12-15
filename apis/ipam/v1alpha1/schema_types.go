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
	Tag                []*NddrIpamIpamNetworkInstanceTag       `json:"tag,omitempty"`
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
	Tag         []*NddrIpamIpamNetworkInstanceIpAddressTag `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpAddressState struct
type NddrIpamIpamNetworkInstanceIpAddressState struct {
	//LastUpdate *string                                              `json:"last-update,omitempty"`
	//Origin     *string                                              `json:"origin,omitempty"`
	IpPrefix []*NddrIpamIpamNetworkInstanceIpAddressStateIpPrefix `json:"ip-prefix,omitempty"`
	IpRange  []*NddrIpamIpamNetworkInstanceIpAddressStateIpRange  `json:"ip-range,omitempty"`
	Reason   *string                                              `json:"reason,omitempty"`
	Status   *string                                              `json:"status,omitempty"`
	Tag      []*NddrIpamIpamNetworkInstanceIpAddressStateTag      `json:"tag,omitempty"`
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

// NddrIpamIpamNetworkInstanceIpAddressStateTag struct
type NddrIpamIpamNetworkInstanceIpAddressStateTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpAddressTag struct
type NddrIpamIpamNetworkInstanceIpAddressTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpPrefix struct
type NddrIpamIpamNetworkInstanceIpPrefix struct {
	AdminState  *string `json:"admin-state,omitempty"`
	Description *string `json:"description,omitempty"`
	Pool        *bool   `json:"pool,omitempty"`
	Prefix      *string `json:"prefix"`
	//RirName     *string                                   `json:"rir-name,omitempty"`
	State *NddrIpamIpamNetworkInstanceIpPrefixState `json:"state,omitempty"`
	Tag   []*NddrIpamIpamNetworkInstanceIpPrefixTag `json:"tag,omitempty"`
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
	Tag    []*NddrIpamIpamNetworkInstanceIpPrefixStateTag  `json:"tag,omitempty"`
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

// NddrIpamIpamNetworkInstanceIpPrefixStateTag struct
type NddrIpamIpamNetworkInstanceIpPrefixStateTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpPrefixTag struct
type NddrIpamIpamNetworkInstanceIpPrefixTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRange struct
type NddrIpamIpamNetworkInstanceIpRange struct {
	AdminState  *string                                  `json:"admin-state,omitempty"`
	Description *string                                  `json:"description,omitempty"`
	End         *string                                  `json:"end"`
	Start       *string                                  `json:"start"`
	State       *NddrIpamIpamNetworkInstanceIpRangeState `json:"state,omitempty"`
	Tag         []*NddrIpamIpamNetworkInstanceIpRangeTag `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRangeState struct
type NddrIpamIpamNetworkInstanceIpRangeState struct {
	//LastUpdate *string                                        `json:"last-update,omitempty"`
	//Origin     *string                                        `json:"origin,omitempty"`
	Parent *NddrIpamIpamNetworkInstanceIpRangeStateParent `json:"parent,omitempty"`
	Reason *string                                        `json:"reason,omitempty"`
	Size   *uint32                                        `json:"size,omitempty"`
	Status *string                                        `json:"status,omitempty"`
	Tag    []*NddrIpamIpamNetworkInstanceIpRangeStateTag  `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRangeStateParent struct
type NddrIpamIpamNetworkInstanceIpRangeStateParent struct {
	IpPrefix []*NddrIpamIpamNetworkInstanceIpRangeStateParentIpPrefix `json:"ip-prefix,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRangeStateParentIpPrefix struct
type NddrIpamIpamNetworkInstanceIpRangeStateParentIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddrIpamIpamNetworkInstanceIpRangeStateTag struct
type NddrIpamIpamNetworkInstanceIpRangeStateTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamNetworkInstanceIpRangeTag struct
type NddrIpamIpamNetworkInstanceIpRangeTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamNetworkInstanceState struct
type NddrIpamIpamNetworkInstanceState struct {
	//LastUpdate *string                                `json:"last-update,omitempty"`
	//Origin     *string                                `json:"origin,omitempty"`
	Reason *string                                `json:"reason,omitempty"`
	Status *string                                `json:"status,omitempty"`
	Tag    []*NddrIpamIpamNetworkInstanceStateTag `json:"tag,omitempty"`
}

// NddrIpamIpamNetworkInstanceStateTag struct
type NddrIpamIpamNetworkInstanceStateTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamNetworkInstanceTag struct
type NddrIpamIpamNetworkInstanceTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamState struct
type NddrIpamIpamState struct {
	//LastUpdate *string                 `json:"last-update,omitempty"`
	//Origin     *string                 `json:"origin,omitempty"`
	Reason *string                 `json:"reason,omitempty"`
	Status *string                 `json:"status,omitempty"`
	Tag    []*NddrIpamIpamStateTag `json:"tag,omitempty"`
}

// NddrIpamIpamStateTag struct
type NddrIpamIpamStateTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddrIpamIpamTag struct
type NddrIpamIpamTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// Root is the root of the schema
type Root struct {
	IpamNddrIpam *NddrIpam `json:"Nddr-ipam,omitempty"`
}
