/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	grafanav1 "github.com/dichque/grafana-operator/pkg/apis/grafana/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGrafanas implements GrafanaInterface
type FakeGrafanas struct {
	Fake *FakeAimsV1
	ns   string
}

var grafanasResource = schema.GroupVersionResource{Group: "aims.cisco.com", Version: "v1", Resource: "grafanas"}

var grafanasKind = schema.GroupVersionKind{Group: "aims.cisco.com", Version: "v1", Kind: "Grafana"}

// Get takes name of the grafana, and returns the corresponding grafana object, and an error if there is any.
func (c *FakeGrafanas) Get(ctx context.Context, name string, options v1.GetOptions) (result *grafanav1.Grafana, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(grafanasResource, c.ns, name), &grafanav1.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*grafanav1.Grafana), err
}

// List takes label and field selectors, and returns the list of Grafanas that match those selectors.
func (c *FakeGrafanas) List(ctx context.Context, opts v1.ListOptions) (result *grafanav1.GrafanaList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(grafanasResource, grafanasKind, c.ns, opts), &grafanav1.GrafanaList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &grafanav1.GrafanaList{ListMeta: obj.(*grafanav1.GrafanaList).ListMeta}
	for _, item := range obj.(*grafanav1.GrafanaList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested grafanas.
func (c *FakeGrafanas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(grafanasResource, c.ns, opts))

}

// Create takes the representation of a grafana and creates it.  Returns the server's representation of the grafana, and an error, if there is any.
func (c *FakeGrafanas) Create(ctx context.Context, grafana *grafanav1.Grafana, opts v1.CreateOptions) (result *grafanav1.Grafana, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(grafanasResource, c.ns, grafana), &grafanav1.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*grafanav1.Grafana), err
}

// Update takes the representation of a grafana and updates it. Returns the server's representation of the grafana, and an error, if there is any.
func (c *FakeGrafanas) Update(ctx context.Context, grafana *grafanav1.Grafana, opts v1.UpdateOptions) (result *grafanav1.Grafana, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(grafanasResource, c.ns, grafana), &grafanav1.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*grafanav1.Grafana), err
}

// Delete takes name of the grafana and deletes it. Returns an error if one occurs.
func (c *FakeGrafanas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(grafanasResource, c.ns, name), &grafanav1.Grafana{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGrafanas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(grafanasResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &grafanav1.GrafanaList{})
	return err
}

// Patch applies the patch and returns the patched grafana.
func (c *FakeGrafanas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *grafanav1.Grafana, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(grafanasResource, c.ns, name, pt, data, subresources...), &grafanav1.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*grafanav1.Grafana), err
}
