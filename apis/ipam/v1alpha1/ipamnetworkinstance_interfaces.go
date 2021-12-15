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

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/utils"
)

var _ In = &IpamNetworkInstance{}

// +k8s:deepcopy-gen=false
type In interface {
	resource.Object
	resource.Conditioned

	GetIpamName() string
	GetAdminState() string
	GetDescription() string
	GetTags() map[string]string
	InitializeResource() error
	SetStatus(string)
	SetReason(string)
	GetStatus() string
}

// GetCondition of this Network Node.
func (x *IpamNetworkInstance) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *IpamNetworkInstance) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *IpamNetworkInstance) GetIpamName() string {
	if reflect.ValueOf(x.Spec.IpamName).IsZero() {
		return ""
	}
	return *x.Spec.IpamName
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
	tags := make([]*NddrIpamIpamNetworkInstanceTag, 0, len(x.Spec.IpamNetworkInstance.Tag))
	for _, tag := range x.Spec.IpamNetworkInstance.Tag {
		tags = append(tags, &NddrIpamIpamNetworkInstanceTag{
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
			Tag:    make([]*NddrIpamIpamNetworkInstanceStateTag, 0),
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
