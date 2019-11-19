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

package integration

import (
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"regexp"
	"testing"

	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
)

const (
	stateValueRegexp     = `192\.[0-9]+\.[0-9]+\.[0-9]+`
	stateControllersPath = "/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/state/address"
)

// TestSingleState tests query of a single GNMI path of a read/only value to a single device
func TestSingleState(t *testing.T) {
	// Get the first configured device from the environment.
	device := env.GetDevices()[0]

	// Make a GNMI client to use for requests
	c, err := env.NewGnmiClient(MakeContext(), "")
	assert.NoError(t, err)
	assert.True(t, c != nil, "Fetching client returned nil")

	// Check that the value was correctly retrieved from the device and store in the state cache
	valueAfter, extensions, errorAfter := GNMIGet(MakeContext(), c, makeDevicePath(device, stateControllersPath))
	assert.NoError(t, errorAfter)
	assert.NotEqual(t, "", valueAfter, "Query after state returned an error: %s\n", errorAfter)
	re := regexp.MustCompile(stateValueRegexp)
	match := re.MatchString(valueAfter[0].pathDataValue)
	assert.Equal(t, match, true, "Query for state returned the wrong value: %s\n", valueAfter)
	assert.Equal(t, 0, len(extensions))
}

func init() {
	test.Registry.RegisterTest("single-state", TestSingleState, []*runner.TestSuite{AllTests, IntegrationTests})
}