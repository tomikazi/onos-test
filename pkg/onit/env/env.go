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
	"github.com/onosproject/onos-test/pkg/kube"
	"github.com/onosproject/onos-test/pkg/onit/cluster"
)

// New returns a new onit Env
func New(kube kube.API) Env {
	return &clusterEnv{
		cluster: cluster.New(kube),
	}
}

// Env is an interface for tests to operate on the ONOS environment
type Env interface {
	// Atomix returns the Atomix environment
	Atomix() Atomix

	// Database returns the database environment
	Database() Database

	// Topo returns the topo environment
	Topo() Topo

	// Config returns the config environment
	Config() Config

	// Simulators returns the simulators environment
	Simulators() Simulators

	// Simulator returns the environment for a simulator by name
	Simulator(name string) Simulator

	// NewSimulator returns a new SimulatorSetup for adding a simulator to the cluster
	NewSimulator() SimulatorSetup

	// Networks returns the networks environment
	Networks() Networks

	// Network returns the environment for a network by name
	Network(name string) Network

	// NewNetwork returns a new NetworkSetup for adding a network to the cluster
	NewNetwork() NetworkSetup

	// Apps returns the applications environment
	Apps() Apps

	// App returns the environment for an app by name
	App(name string) App

	// NewApp returns a new AppSetup for adding an application to the cluster
	NewApp() AppSetup
}

// clusterEnv is an implementation of the Env interface
type clusterEnv struct {
	cluster *cluster.Cluster
}

func (e *clusterEnv) Atomix() Atomix {
	return &clusterAtomix{
		clusterService: &clusterService{
			service: e.cluster.Atomix().Service,
		},
	}
}

func (e *clusterEnv) Database() Database {
	return &clusterDatabase{
		database: e.cluster.Database(),
	}
}

func (e *clusterEnv) Topo() Topo {
	return &clusterTopo{
		clusterService: &clusterService{
			service: e.cluster.Topo().Service,
		},
	}
}

func (e *clusterEnv) Config() Config {
	return &clusterConfig{
		clusterService: &clusterService{
			service: e.cluster.Config().Service,
		},
	}
}

func (e *clusterEnv) Simulators() Simulators {
	return &clusterSimulators{
		simulators: e.cluster.Simulators(),
	}
}

func (e *clusterEnv) Simulator(name string) Simulator {
	return e.Simulators().Get(name)
}

func (e *clusterEnv) NewSimulator() SimulatorSetup {
	return e.Simulators().New()
}

func (e *clusterEnv) Networks() Networks {
	return &clusterNetworks{
		networks: e.cluster.Networks(),
	}
}

func (e *clusterEnv) Network(name string) Network {
	return e.Networks().Get(name)
}

func (e *clusterEnv) NewNetwork() NetworkSetup {
	return e.Networks().New()
}

func (e *clusterEnv) Apps() Apps {
	return &clusterApps{
		apps: e.cluster.Apps(),
	}
}

func (e *clusterEnv) App(name string) App {
	return e.Apps().Get(name)
}

func (e *clusterEnv) NewApp() AppSetup {
	return e.Apps().New()
}