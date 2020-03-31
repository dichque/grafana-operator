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

package v1

import (
	"time"

	v1 "github.com/dichque/grafana-operator/pkg/apis/grafana/v1"
	scheme "github.com/dichque/grafana-operator/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// GrafanasGetter has a method to return a GrafanaInterface.
// A group's client should implement this interface.
type GrafanasGetter interface {
	Grafanas(namespace string) GrafanaInterface
}

// GrafanaInterface has methods to work with Grafana resources.
type GrafanaInterface interface {
	Create(*v1.Grafana) (*v1.Grafana, error)
	Update(*v1.Grafana) (*v1.Grafana, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Grafana, error)
	List(opts metav1.ListOptions) (*v1.GrafanaList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Grafana, err error)
	GrafanaExpansion
}

// grafanas implements GrafanaInterface
type grafanas struct {
	client rest.Interface
	ns     string
}

// newGrafanas returns a Grafanas
func newGrafanas(c *AimsV1Client, namespace string) *grafanas {
	return &grafanas{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the grafana, and returns the corresponding grafana object, and an error if there is any.
func (c *grafanas) Get(name string, options metav1.GetOptions) (result *v1.Grafana, err error) {
	result = &v1.Grafana{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("grafanas").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Grafanas that match those selectors.
func (c *grafanas) List(opts metav1.ListOptions) (result *v1.GrafanaList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.GrafanaList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("grafanas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested grafanas.
func (c *grafanas) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("grafanas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a grafana and creates it.  Returns the server's representation of the grafana, and an error, if there is any.
func (c *grafanas) Create(grafana *v1.Grafana) (result *v1.Grafana, err error) {
	result = &v1.Grafana{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("grafanas").
		Body(grafana).
		Do().
		Into(result)
	return
}

// Update takes the representation of a grafana and updates it. Returns the server's representation of the grafana, and an error, if there is any.
func (c *grafanas) Update(grafana *v1.Grafana) (result *v1.Grafana, err error) {
	result = &v1.Grafana{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("grafanas").
		Name(grafana.Name).
		Body(grafana).
		Do().
		Into(result)
	return
}

// Delete takes name of the grafana and deletes it. Returns an error if one occurs.
func (c *grafanas) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("grafanas").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *grafanas) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("grafanas").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched grafana.
func (c *grafanas) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Grafana, err error) {
	result = &v1.Grafana{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("grafanas").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
