/*
Copyright 2021 NDD.

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
	"strings"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/utils"
	nddov1 "github.com/yndd/nddo-runtime/apis/common/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ InList = &IpamNetworkInstanceList{}

// +k8s:deepcopy-gen=false
type InList interface {
	client.ObjectList

	GetNetworkInstances() []In
}

func (x *IpamNetworkInstanceList) GetNetworkInstances() []In {
	xs := make([]In, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ In = &IpamNetworkInstance{}

// +k8s:deepcopy-gen=false
type In interface {
	resource.Object
	resource.Conditioned

	GetOrganizationName() string
	GetIpamName() string
	GetNetworkInstanceName() string
	GetAdminState() string
	GetDescription() string
	GetAllocationStrategy() string
	GetDefaultPrefixLength(string, string) *uint32
	GetTags() map[string]string
	InitializeResource() error
	SetStatus(string)
	SetReason(string)
	GetStatus() string

	SetOrganizationName(string)
	SetIpamName(string)
	SetNetworkInstanceName(string)
}

// GetCondition of this Network Node.
func (x *IpamNetworkInstance) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *IpamNetworkInstance) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *IpamNetworkInstance) GetOrganizationName() string {
	split := strings.Split(x.GetName(), ".")
	if len(split) >= 3 {
		return split[0]
	}
	return ""
}

func (x *IpamNetworkInstance) GetIpamName() string {
	split := strings.Split(x.GetName(), ".")
	if len(split) >= 3 {
		return split[1]
	}
	return ""
}

func (x *IpamNetworkInstance) GetNetworkInstanceName() string {
	split := strings.Split(x.GetName(), ".")
	if len(split) >= 3 {
		return split[2]
	}
	return ""
}

func (x *IpamNetworkInstance) GetAdminState() string {
	if reflect.ValueOf(x.Spec.IpamNetworkInstance.AdminState).IsZero() {
		return ""
	}
	return *x.Spec.IpamNetworkInstance.AdminState
}

func (x *IpamNetworkInstance) GetDescription() string {
	if reflect.ValueOf(x.Spec.IpamNetworkInstance.Description).IsZero() {
		return ""
	}
	return *x.Spec.IpamNetworkInstance.Description
}

func (x *IpamNetworkInstance) GetAllocationStrategy() string {
	if reflect.ValueOf(x.Spec.IpamNetworkInstance.AllocationStrategy).IsZero() {
		return ""
	}
	return *x.Spec.IpamNetworkInstance.AllocationStrategy
}

func (x *IpamNetworkInstance) GetDefaultPrefixLength(p, af string) *uint32 {
	if reflect.ValueOf(x.Spec.IpamNetworkInstance.DefaultPrefixLength).IsZero() {
		return nil
	}
	if purpose, ok := x.Spec.IpamNetworkInstance.DefaultPrefixLength[p]; ok {
		if pl, ok := purpose.AddressFamily[af]; ok {
			return pl
		}
	}
	return nil
}

func (x *IpamNetworkInstance) GetTags() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.IpamNetworkInstance.Tag).IsZero() {
		return s
	}
	for _, tag := range x.Spec.IpamNetworkInstance.Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *IpamNetworkInstance) InitializeResource() error {
	tags := make([]*nddov1.Tag, 0, len(x.Spec.IpamNetworkInstance.Tag))
	for _, tag := range x.Spec.IpamNetworkInstance.Tag {
		tags = append(tags, &nddov1.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	if x.Status.IpamNetworkInstance != nil {
		// pool was already initialiazed
		// copy the spec, but not the state
		x.Status.IpamNetworkInstance.AdminState = x.Spec.IpamNetworkInstance.AdminState
		x.Status.IpamNetworkInstance.Description = x.Spec.IpamNetworkInstance.Description
		x.Status.IpamNetworkInstance.AllocationStrategy = x.Spec.IpamNetworkInstance.AllocationStrategy
		x.Status.IpamNetworkInstance.Tag = tags
		return nil
	}

	x.Status.IpamNetworkInstance = &NddrIpamIpamNetworkInstance{
		AdminState:         x.Spec.IpamNetworkInstance.AdminState,
		Description:        x.Spec.IpamNetworkInstance.Description,
		AllocationStrategy: x.Spec.IpamNetworkInstance.AllocationStrategy,
		Tag:                tags,
		State: &NddrIpamIpamNetworkInstanceState{
			Status: utils.StringPtr(""),
			Reason: utils.StringPtr(""),
			Tag:    make([]*nddov1.Tag, 0),
		},
	}
	return nil
}

func (x *IpamNetworkInstance) SetStatus(s string) {
	x.Status.IpamNetworkInstance.State.Status = &s
}

func (x *IpamNetworkInstance) SetReason(s string) {
	x.Status.IpamNetworkInstance.State.Reason = &s
}

func (x *IpamNetworkInstance) GetStatus() string {
	if x.Status.IpamNetworkInstance != nil && x.Status.IpamNetworkInstance.State != nil && x.Status.IpamNetworkInstance.State.Status != nil {
		return *x.Status.IpamNetworkInstance.State.Status
	}
	return "unknown"
}

func (x *IpamNetworkInstance) SetOrganizationName(s string) {
	x.Status.OrganizationName = &s
}

func (x *IpamNetworkInstance) SetIpamName(s string) {
	x.Status.IpamName = &s
}

func (x *IpamNetworkInstance) SetNetworkInstanceName(s string) {
	x.Status.NetworkInstanceName = &s
}
