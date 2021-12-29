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

var _ IpList = &IpamList{}

// +k8s:deepcopy-gen=false
type IpList interface {
	client.ObjectList

	GetIpams() []Ip
}

func (x *IpamList) GetIpams() []Ip {
	xs := make([]Ip, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ Ip = &Ipam{}

// +k8s:deepcopy-gen=false
type Ip interface {
	resource.Object
	resource.Conditioned

	GetOrganizationName() string
	GetIpamName() string
	GetAdminState() string
	GetDescription() string
	InitializeResource() error
	SetStatus(string)
	SetReason(string)
	GetStatus() string

	SetOrganizationName(string)
	SetIpamName(string)
}

// GetCondition of this Network Node.
func (x *Ipam) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *Ipam) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *Ipam) GetOrganizationName() string {
	split := strings.Split(x.GetName(), ".")
	if len(split) >= 2 {
		return split[0]
	}
	return ""
}

func (x *Ipam) GetIpamName() string {
	split := strings.Split(x.GetName(), ".")
	if len(split) >= 2 {
		return split[1]
	}
	return ""
}

func (x *Ipam) GetAdminState() string {
	if reflect.ValueOf(x.Spec.Ipam.AdminState).IsZero() {
		return ""
	}
	return *x.Spec.Ipam.AdminState
}

func (x *Ipam) GetDescription() string {
	if reflect.ValueOf(x.Spec.Ipam.Description).IsZero() {
		return ""
	}
	return *x.Spec.Ipam.Description
}

func (x *Ipam) InitializeResource() error {
	if x.Status.Ipam != nil {
		// pool was already initialiazed
		// copy the spec, but not the state
		x.Status.Ipam.AdminState = x.Spec.Ipam.AdminState
		x.Status.Ipam.Description = x.Spec.Ipam.Description
		return nil
	}

	x.Status.Ipam = &NddrIpamIpam{
		AdminState:  x.Spec.Ipam.AdminState,
		Description: x.Spec.Ipam.Description,
		State: &NddrIpamIpamState{
			Status: utils.StringPtr(""),
			Reason: utils.StringPtr(""),
			Tag:    make([]*nddov1.Tag, 0),
		},
	}
	return nil
}

func (x *Ipam) SetStatus(s string) {
	x.Status.Ipam.State.Status = &s
}

func (x *Ipam) SetReason(s string) {
	x.Status.Ipam.State.Reason = &s
}

func (x *Ipam) GetStatus() string {
	if x.Status.Ipam != nil && x.Status.Ipam.State != nil && x.Status.Ipam.State.Status != nil {
		return *x.Status.Ipam.State.Status
	}
	return "unknown"
}

func (x *Ipam) SetOrganizationName(s string) {
	x.Status.OrganizationName = &s
}

func (x *Ipam) SetIpamName(s string) {
	x.Status.IpamName = &s
}
