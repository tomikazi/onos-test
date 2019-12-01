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

package onit

import "github.com/onosproject/onos-test/pkg/kubetest"

// RegisterTests registers a test suite
func RegisterTests(name string, suite kubetest.TestingSuite) {
	kubetest.RegisterTests(name, suite)
}

// RegisterBenchmarks registers a benchmark suite
func RegisterBenchmarks(name string, suite kubetest.BenchmarkingSuite) {
	kubetest.RegisterBenchmarks(name, suite)
}

// RegisterScripts registers a script suite
func RegisterScripts(name string, suite kubetest.ScriptingSuite) {
	kubetest.RegisterScripts(name, suite)
}
