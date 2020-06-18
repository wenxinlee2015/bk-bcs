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

package v1alpha1

import (
	v1alpha1 "bk-bcs/bcs-k8s/bcs-gamestatefulset-operator/pkg/apis/tkex/v1alpha1"
	scheme "bk-bcs/bcs-k8s/bcs-gamestatefulset-operator/pkg/clientset/internalclientset/scheme"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	autoscaling "k8s.io/kubernetes/pkg/apis/autoscaling"
)

// GameStatefulSetsGetter has a method to return a GameStatefulSetInterface.
// A group's client should implement this interface.
type GameStatefulSetsGetter interface {
	GameStatefulSets(namespace string) GameStatefulSetInterface
}

// GameStatefulSetInterface has methods to work with GameStatefulSet resources.
type GameStatefulSetInterface interface {
	Create(*v1alpha1.GameStatefulSet) (*v1alpha1.GameStatefulSet, error)
	Update(*v1alpha1.GameStatefulSet) (*v1alpha1.GameStatefulSet, error)
	UpdateStatus(*v1alpha1.GameStatefulSet) (*v1alpha1.GameStatefulSet, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.GameStatefulSet, error)
	List(opts v1.ListOptions) (*v1alpha1.GameStatefulSetList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GameStatefulSet, err error)
	GetScale(gameStatefulSetName string, options v1.GetOptions) (*autoscaling.Scale, error)
	UpdateScale(gameStatefulSetName string, scale *autoscaling.Scale) (*autoscaling.Scale, error)

	GameStatefulSetExpansion
}

// gameStatefulSets implements GameStatefulSetInterface
type gameStatefulSets struct {
	client rest.Interface
	ns     string
}

// newGameStatefulSets returns a GameStatefulSets
func newGameStatefulSets(c *TkexV1alpha1Client, namespace string) *gameStatefulSets {
	return &gameStatefulSets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the gameStatefulSet, and returns the corresponding gameStatefulSet object, and an error if there is any.
func (c *gameStatefulSets) Get(name string, options v1.GetOptions) (result *v1alpha1.GameStatefulSet, err error) {
	result = &v1alpha1.GameStatefulSet{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of GameStatefulSets that match those selectors.
func (c *gameStatefulSets) List(opts v1.ListOptions) (result *v1alpha1.GameStatefulSetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.GameStatefulSetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested gameStatefulSets.
func (c *gameStatefulSets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a gameStatefulSet and creates it.  Returns the server's representation of the gameStatefulSet, and an error, if there is any.
func (c *gameStatefulSets) Create(gameStatefulSet *v1alpha1.GameStatefulSet) (result *v1alpha1.GameStatefulSet, err error) {
	result = &v1alpha1.GameStatefulSet{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		Body(gameStatefulSet).
		Do().
		Into(result)
	return
}

// Update takes the representation of a gameStatefulSet and updates it. Returns the server's representation of the gameStatefulSet, and an error, if there is any.
func (c *gameStatefulSets) Update(gameStatefulSet *v1alpha1.GameStatefulSet) (result *v1alpha1.GameStatefulSet, err error) {
	result = &v1alpha1.GameStatefulSet{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		Name(gameStatefulSet.Name).
		Body(gameStatefulSet).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *gameStatefulSets) UpdateStatus(gameStatefulSet *v1alpha1.GameStatefulSet) (result *v1alpha1.GameStatefulSet, err error) {
	result = &v1alpha1.GameStatefulSet{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		Name(gameStatefulSet.Name).
		SubResource("status").
		Body(gameStatefulSet).
		Do().
		Into(result)
	return
}

// Delete takes name of the gameStatefulSet and deletes it. Returns an error if one occurs.
func (c *gameStatefulSets) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *gameStatefulSets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched gameStatefulSet.
func (c *gameStatefulSets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GameStatefulSet, err error) {
	result = &v1alpha1.GameStatefulSet{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("gamestatefulsets").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}

// GetScale takes name of the gameStatefulSet, and returns the corresponding autoscaling.Scale object, and an error if there is any.
func (c *gameStatefulSets) GetScale(gameStatefulSetName string, options v1.GetOptions) (result *autoscaling.Scale, err error) {
	result = &autoscaling.Scale{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		Name(gameStatefulSetName).
		SubResource("scale").
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// UpdateScale takes the top resource name and the representation of a scale and updates it. Returns the server's representation of the scale, and an error, if there is any.
func (c *gameStatefulSets) UpdateScale(gameStatefulSetName string, scale *autoscaling.Scale) (result *autoscaling.Scale, err error) {
	result = &autoscaling.Scale{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("gamestatefulsets").
		Name(gameStatefulSetName).
		SubResource("scale").
		Body(scale).
		Do().
		Into(result)
	return
}
