/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
// Code generated by client-gen. DO NOT EDIT.

package v2

import (
	v2 "github.com/Tencent/bk-bcs/bcs-mesos/mesosv2/apis/bkbcs/v2"
	"github.com/Tencent/bk-bcs/bcs-mesos/mesosv2/generated/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type BkbcsV2Interface interface {
	RESTClient() rest.Interface
	AgentsGetter
	TaskGroupsGetter
}

// BkbcsV2Client is used to interact with features provided by the bkbcs group.
type BkbcsV2Client struct {
	restClient rest.Interface
}

func (c *BkbcsV2Client) Agents(namespace string) AgentInterface {
	return newAgents(c, namespace)
}

func (c *BkbcsV2Client) TaskGroups(namespace string) TaskGroupInterface {
	return newTaskGroups(c, namespace)
}

// NewForConfig creates a new BkbcsV2Client for the given config.
func NewForConfig(c *rest.Config) (*BkbcsV2Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &BkbcsV2Client{client}, nil
}

// NewForConfigOrDie creates a new BkbcsV2Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *BkbcsV2Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new BkbcsV2Client for the given RESTClient.
func New(c rest.Interface) *BkbcsV2Client {
	return &BkbcsV2Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v2.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *BkbcsV2Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}