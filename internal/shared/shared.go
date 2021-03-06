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
package shared

import (
	"time"

	"github.com/hansthienpondt/goipam/pkg/table"
	"github.com/yndd/ndd-runtime/pkg/logging"
)

type NddControllerOptions struct {
	Logger    logging.Logger
	Poll      time.Duration
	Namespace string
	// first map represents namespace/organization/ipam, second map represents network-instance
	Iptree map[string]*table.RouteTable
	//Yentry      *yentry.Entry
	//GnmiAddress string
}
