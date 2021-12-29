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

	"github.com/hansthienpondt/goipam/pkg/table"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ipamv1alpha1 "github.com/yndd/nddr-ipam/apis/ipam/v1alpha1"
)

type Config struct {
	// Address
	Address string
	// Generic
	MaxSubscriptions int64
	MaxUnaryRPC      int64
	// TLS
	InSecure   bool
	SkipVerify bool
	CaFile     string
	CertFile   string
	KeyFile    string
	// observability
	EnableMetrics bool
	Debug         bool
}

// Option can be used to manipulate Options.
type Option func(Server)

// WithLogger specifies how the Reconciler should log messages.
func WithLogger(log logging.Logger) Option {
	return func(s Server) {
		s.WithLogger(log)
	}
}

func WithConfig(cfg Config) Option {
	return func(s Server) {
		s.WithConfig(cfg)
	}
}

func WithIpTree(iptree map[string]*table.RouteTable) Option {
	return func(s Server) {
		s.WithIpTree(iptree)
	}
}

func WithClient(c client.Client) Option {
	return func(s Server) {
		s.WithClient(c)
	}
}

func WithNewResourceFn(f func() ipamv1alpha1.In) Option {
	return func(r Server) {
		r.WithNewResourceFn(f)
	}
}

type Server interface {
	WithLogger(log logging.Logger)
	WithConfig(cfg Config)
	WithIpTree(iptree map[string]*table.RouteTable)
	WithClient(a client.Client)
	WithNewResourceFn(f func() ipamv1alpha1.In)
	Run(ctx context.Context) error
}
