// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package env

import (
	"context"
	"github.com/onosproject/onos-config/pkg/northbound/admin"
	"github.com/openconfig/gnmi/client"
	gnmi "github.com/openconfig/gnmi/client/gnmi"
	"time"
)

// Config provides the config environment
type Config interface {
	Service

	// Destination returns the gNMI client destination
	Destination() client.Destination

	// NewAdminServiceClient returns the config AdminService client
	NewAdminServiceClient() (admin.ConfigAdminServiceClient, error)

	// NewGNMIClient returns the gNMI client
	NewGNMIClient() (*gnmi.Client, error)
}

var _ Config = &clusterConfig{}

// clusterConfig is an implementation of the Config interface
type clusterConfig struct {
	*clusterService
}

func (e *clusterConfig) Destination() client.Destination {
	return client.Destination{
		Addrs:   []string{e.Address()},
		Target:  "gnmi",
		TLS:     e.Credentials(),
		Timeout: 10 * time.Second,
	}
}

func (e *clusterConfig) NewAdminServiceClient() (admin.ConfigAdminServiceClient, error) {
	conn, err := e.Connect()
	if err != nil {
		return nil, err
	}
	return admin.NewConfigAdminServiceClient(conn), nil
}

func (e *clusterConfig) NewGNMIClient() (*gnmi.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := gnmi.New(ctx, e.Destination())
	if err != nil {
		return nil, err
	}
	return client.(*gnmi.Client), nil
}