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

package grpcserver

import (
	"context"
	"strings"
	"time"

	"github.com/hansthienpondt/goipam/pkg/table"
	"github.com/pkg/errors"
	"github.com/yndd/nddo-grpc/resource/resourcepb"
	ipamv1alpha1 "github.com/yndd/nddr-ipam/apis/ipam/v1alpha1"
	"inet.af/netaddr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
)

func (s *server) ResourceGet(ctx context.Context, req *resourcepb.Request) (*resourcepb.Reply, error) {
	log := s.log.WithValues("Request", req)
	log.Debug("ResourceGet...")

	return &resourcepb.Reply{Ready: true}, nil
}

func (s *server) ResourceAlloc(ctx context.Context, req *resourcepb.Request) (*resourcepb.Reply, error) {
	log := s.log.WithValues("Request", req)
	log.Debug("ResourceAlloc...")

	namespace := req.GetNamespace()
	niName := req.GetResourceName()
	crName := strings.Join([]string{namespace, niName}, ".")

	ni := &ipamv1alpha1.IpamNetworkInstance{}
	if err := s.client.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      niName}, ni); err != nil {
		// can happen when the networkinstance is not found
		log.Debug("NetworkInstance not available")
		return &resourcepb.Reply{Ready: false}, errors.New("NetworkInstance not available")
	}

	if ni.GetCondition(ipamv1alpha1.ConditionKindReady).Status != corev1.ConditionTrue {
		log.Debug("NetworkInstance not ready")
		return &resourcepb.Reply{Ready: false}, errors.New("NetworkInstance not ready")
	}

	if _, ok := s.iptree[crName]; !ok {
		log.Debug("Routing table not ready")
		return &resourcepb.Reply{Ready: false}, nil
	}

	// the selector is used in the tree to find the entry in the tree
	// we use all the keys in the source-tag and selector for the search
	fullselector := labels.NewSelector()
	l := make(map[string]string)
	for key, val := range req.GetAlloc().GetSelector() {
		req, err := labels.NewRequirement(key, selection.In, []string{val})
		if err != nil {
			log.Debug("wrong object", "Error", err)
			return &resourcepb.Reply{Ready: true}, err
		}
		fullselector = fullselector.Add(*req)
		l[key] = val
	}
	for key, val := range req.GetAlloc().GetSourceTag() {
		req, err := labels.NewRequirement(key, selection.In, []string{val})
		if err != nil {
			log.Debug("wrong object", "Error", err)
			return &resourcepb.Reply{Ready: true}, err
		}
		fullselector = fullselector.Add(*req)
		l[key] = val
	}
	prefix := req.GetAlloc().GetIpPrefix()
	if prefix != "" {
		log.Debug("alloc has prefix", "prefix", prefix)

		// via selector perform allocation
		selector := labels.NewSelector()
		for key, val := range req.GetAlloc().GetSelector() {
			req, err := labels.NewRequirement(key, selection.In, []string{val})
			if err != nil {
				log.Debug("wrong object", "Error", err)
				return &resourcepb.Reply{Ready: true}, err
			}
			selector = selector.Add(*req)
		}

		a, err := netaddr.ParseIPPrefix(prefix)
		if err != nil {
			log.Debug("Cannot parse ip prefix", "error", err)
			return &resourcepb.Reply{Ready: true}, errors.Wrap(err, "Cannot parse ip prefix")
		}
		route := table.NewRoute(a)
		route.UpdateLabel(l)

		if err := s.iptree[crName].Add(route); err != nil {
			log.Debug("route insertion failed")
			if !strings.Contains(err.Error(), "already exists") {
				return &resourcepb.Reply{Ready: true}, errors.Wrap(err, "route insertion failed")
			}
		}
		prefix = route.String()
	} else {
		// alloc has no prefix assigned, try to assign prefix
		// check if the prefix already exists
		routes := s.iptree[crName].GetByLabel(fullselector)
		if len(routes) == 0 {
			// allocate prefix
			log.Debug("Query not found, allocate a prefix")

			// via selector perform allocation
			selector := labels.NewSelector()
			for key, val := range req.GetAlloc().GetSelector() {
				req, err := labels.NewRequirement(key, selection.In, []string{val})
				if err != nil {
					log.Debug("wrong object", "Error", err)
					return &resourcepb.Reply{Ready: true}, errors.Wrap(err, "wrong object")
				}
				selector = selector.Add(*req)
			}

			routes := s.iptree[crName].GetByLabel(selector)

			// required during startup when not everything is initialized
			// we break and the reconciliation will take care
			if len(routes) == 0 {
				log.Debug("no available routes")
				return &resourcepb.Reply{Ready: true}, errors.New("no available routes")
			}

			// TBD we take the first prefix
			prefixLength, err := getPrefixLength(req, ni)
			if err != nil {
				return &resourcepb.Reply{Ready: true}, errors.Wrap(err, "prefix Length not properly configured")
			}

			a, ok := s.iptree[crName].FindFreePrefix(routes[0].IPPrefix(), uint8(prefixLength))
			if !ok {
				log.Debug("allocation failed")
				return &resourcepb.Reply{Ready: true}, errors.New("allocation failed")
			}

			route := table.NewRoute(a)
			route.UpdateLabel(l)
			if err := s.iptree[crName].Add(route); err != nil {
				log.Debug("route insertion failed")
				return &resourcepb.Reply{Ready: true}, errors.Wrap(err, "route insertion failed")
			}
			prefix = route.String()

		} else {
			if len(routes) > 1 {
				// this should never happen since the labels should provide uniqueness
				log.Debug("strange situation, route in tree found multiple times", "ases", routes)
			}
			prefix = routes[0].IPPrefix().String()
		}

	}

	return &resourcepb.Reply{
		Ready:      true,
		Timestamp:  time.Now().UnixNano(),
		ExpiryTime: time.Now().UnixNano(),
		Data: map[string]*resourcepb.TypedValue{
			"ip-prefix": {Value: &resourcepb.TypedValue_StringVal{StringVal: prefix}},
		},
	}, nil
}

func (s *server) ResourceDeAlloc(ctx context.Context, req *resourcepb.Request) (*resourcepb.Reply, error) {
	log := s.log.WithValues("Request", req)
	log.Debug("ResourceDeAlloc...")

	return &resourcepb.Reply{Ready: true}, nil
}

func getPrefixLength(req *resourcepb.Request, ni *ipamv1alpha1.IpamNetworkInstance) (uint32, error) {
	//prefixLength := cr.GetPrefixLength()
	//if prefixLength == nil {
	purpose := req.GetAlloc().GetSelector()[ipamv1alpha1.KeyPurpose]
	af := req.GetAlloc().GetSelector()[ipamv1alpha1.KeyAddressFamily]
	prefixLength := ni.GetDefaultPrefixLength(purpose, af)
	if prefixLength == nil {
		return 0, errors.New("default prefix length not configured properly")
	}
	//}
	return *prefixLength, nil
}
