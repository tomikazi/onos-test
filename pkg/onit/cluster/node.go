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

package cluster

import (
	"crypto/tls"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func newNode(name string, port int, image string, client *client) *Node {
	return &Node{
		client: client,
		name:   name,
		port:   port,
		image:  image,
	}
}

// Node provides the environment for a single node
type Node struct {
	*client
	name       string
	port       int
	image      string
	pullPolicy corev1.PullPolicy
}

// Name returns the node name
func (n *Node) Name() string {
	return n.name
}

// SetName sets the node name
func (n *Node) SetName(name string) {
	n.name = name
}

// Address returns the service address
func (n *Node) Address() string {
	return fmt.Sprintf("%s:%d", n.name, n.port)
}

// Image returns the image configured for the node
func (n *Node) Image() string {
	return n.image
}

// SetImage sets the node image
func (n *Node) SetImage(image string) {
	n.image = image
}

// PullPolicy returns the image pull policy configured for the node
func (n *Node) PullPolicy() corev1.PullPolicy {
	return n.pullPolicy
}

// SetPullPolicy sets the image pull policy for the node
func (n *Node) SetPullPolicy(pullPolicy corev1.PullPolicy) {
	n.pullPolicy = pullPolicy
}

// Credentials returns the TLS credentials
func (n *Node) Credentials() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

// Connect creates a gRPC client connection to the node
func (n *Node) Connect() (*grpc.ClientConn, error) {
	tlsConfig, err := n.Credentials()
	if err != nil {
		return nil, err
	}
	return grpc.Dial(n.Address(), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}

// AwaitReady waits for the node to become ready
func (n *Node) AwaitReady() error {
	for {
		ready, err := n.isReady()
		if err != nil {
			return err
		} else if ready {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// isReady returns a bool indicating whether the node is ready
func (n *Node) isReady() (bool, error) {
	pod, err := n.kubeClient.CoreV1().Pods(n.namespace).Get(n.name, metav1.GetOptions{})
	if err != nil {
		return false, err
	} else if pod == nil {
		return false, errors.New("node not found")
	}

	for _, status := range pod.Status.ContainerStatuses {
		if !status.Ready {
			return false, nil
		}
	}
	return true, nil
}

// Delete deletes the node
func (n *Node) Delete() error {
	return n.kubeClient.CoreV1().Pods(n.namespace).Delete(n.name, &metav1.DeleteOptions{})
}