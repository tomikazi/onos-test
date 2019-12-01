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

package cli

import (
	"github.com/onosproject/onos-test/pkg/kubetest"
	"github.com/onosproject/onos-test/pkg/util/random"
	corev1 "k8s.io/api/core/v1"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	runExample = `
		# Run a single test on the cluster
		onit run test <name of a test>

		# Run a benchmark on the cluster
		onit run bench <name of a benchmark>`
)

func getRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run an image on Kubernetes",
		RunE:  runRunCommand,
	}
	cmd.Flags().StringP("image", "i", "", "the script image to run")
	cmd.Flags().String("image-pull-policy", string(corev1.PullIfNotPresent), "the Docker image pull policy")
	cmd.Flags().StringP("suite", "s", "", "the script suite to run")
	cmd.Flags().StringP("script", "t", "", "the name of the script method to run")
	cmd.Flags().Duration("timeout", 10*time.Minute, "benchmark timeout")
	return cmd
}

func runRunCommand(cmd *cobra.Command, _ []string) error {
	runCommand(cmd)

	clusterID, _ := cmd.Flags().GetString("cluster")
	image, _ := cmd.Flags().GetString("image")
	suite, _ := cmd.Flags().GetString("suite")
	script, _ := cmd.Flags().GetString("script")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	noTeardown, _ := cmd.Flags().GetBool("no-teardown")
	imagePullPolicy, _ := cmd.Flags().GetString("image-pull-policy")
	pullPolicy := corev1.PullPolicy(imagePullPolicy)

	config := &kubetest.TestConfig{
		TestID:     random.NewPetName(2),
		Type:       kubetest.TestTypeScript,
		Image:      image,
		Suite:      suite,
		Test:       script,
		Timeout:    timeout,
		PullPolicy: pullPolicy,
		Teardown:   !noTeardown,
	}

	// If the cluster ID was not specified, create a new cluster to run the test
	// Otherwise, deploy the test in the existing cluster
	if clusterID == "" {
		runner, err := kubetest.NewTestRunner(config)
		if err != nil {
			return err
		}
		return runner.Run()
	}

	cluster, err := kubetest.NewTestCluster(clusterID)
	if err != nil {
		return err
	}
	if err := cluster.StartTest(config); err != nil {
		return err
	}
	if err := cluster.AwaitTestComplete(config); err != nil {
		return err
	}
	_, code, err := cluster.GetTestResult(config)
	if err != nil {
		return err
	}
	os.Exit(code)
	return nil
}
