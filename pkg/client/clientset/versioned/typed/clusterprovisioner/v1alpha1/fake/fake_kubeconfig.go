/*
Copyright 2018 The Kubernetes Authors.

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

package fake

import (
	v1alpha1 "github.com/tan208123/navigate/pkg/apis/clusterprovisioner/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKubeconfigs implements KubeconfigInterface
type FakeKubeconfigs struct {
	Fake *FakeClusterprovisionerV1alpha1
}

var kubeconfigsResource = schema.GroupVersionResource{Group: "clusterprovisioner.rke.io", Version: "v1alpha1", Resource: "kubeconfigs"}

var kubeconfigsKind = schema.GroupVersionKind{Group: "clusterprovisioner.rke.io", Version: "v1alpha1", Kind: "Kubeconfig"}

// Get takes name of the kubeconfig, and returns the corresponding kubeconfig object, and an error if there is any.
func (c *FakeKubeconfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.Kubeconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(kubeconfigsResource, name), &v1alpha1.Kubeconfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kubeconfig), err
}

// List takes label and field selectors, and returns the list of Kubeconfigs that match those selectors.
func (c *FakeKubeconfigs) List(opts v1.ListOptions) (result *v1alpha1.KubeconfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(kubeconfigsResource, kubeconfigsKind, opts), &v1alpha1.KubeconfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.KubeconfigList{}
	for _, item := range obj.(*v1alpha1.KubeconfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kubeconfigs.
func (c *FakeKubeconfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(kubeconfigsResource, opts))
}

// Create takes the representation of a kubeconfig and creates it.  Returns the server's representation of the kubeconfig, and an error, if there is any.
func (c *FakeKubeconfigs) Create(kubeconfig *v1alpha1.Kubeconfig) (result *v1alpha1.Kubeconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(kubeconfigsResource, kubeconfig), &v1alpha1.Kubeconfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kubeconfig), err
}

// Update takes the representation of a kubeconfig and updates it. Returns the server's representation of the kubeconfig, and an error, if there is any.
func (c *FakeKubeconfigs) Update(kubeconfig *v1alpha1.Kubeconfig) (result *v1alpha1.Kubeconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(kubeconfigsResource, kubeconfig), &v1alpha1.Kubeconfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kubeconfig), err
}

// Delete takes name of the kubeconfig and deletes it. Returns an error if one occurs.
func (c *FakeKubeconfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(kubeconfigsResource, name), &v1alpha1.Kubeconfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKubeconfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(kubeconfigsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.KubeconfigList{})
	return err
}

// Patch applies the patch and returns the patched kubeconfig.
func (c *FakeKubeconfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Kubeconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(kubeconfigsResource, name, data, subresources...), &v1alpha1.Kubeconfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kubeconfig), err
}
