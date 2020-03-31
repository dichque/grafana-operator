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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/dichque/grafana-operator/pkg/apis/grafana/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GrafanaLister helps list Grafanas.
type GrafanaLister interface {
	// List lists all Grafanas in the indexer.
	List(selector labels.Selector) (ret []*v1.Grafana, err error)
	// Grafanas returns an object that can list and get Grafanas.
	Grafanas(namespace string) GrafanaNamespaceLister
	GrafanaListerExpansion
}

// grafanaLister implements the GrafanaLister interface.
type grafanaLister struct {
	indexer cache.Indexer
}

// NewGrafanaLister returns a new GrafanaLister.
func NewGrafanaLister(indexer cache.Indexer) GrafanaLister {
	return &grafanaLister{indexer: indexer}
}

// List lists all Grafanas in the indexer.
func (s *grafanaLister) List(selector labels.Selector) (ret []*v1.Grafana, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Grafana))
	})
	return ret, err
}

// Grafanas returns an object that can list and get Grafanas.
func (s *grafanaLister) Grafanas(namespace string) GrafanaNamespaceLister {
	return grafanaNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GrafanaNamespaceLister helps list and get Grafanas.
type GrafanaNamespaceLister interface {
	// List lists all Grafanas in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Grafana, err error)
	// Get retrieves the Grafana from the indexer for a given namespace and name.
	Get(name string) (*v1.Grafana, error)
	GrafanaNamespaceListerExpansion
}

// grafanaNamespaceLister implements the GrafanaNamespaceLister
// interface.
type grafanaNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Grafanas in the indexer for a given namespace.
func (s grafanaNamespaceLister) List(selector labels.Selector) (ret []*v1.Grafana, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Grafana))
	})
	return ret, err
}

// Get retrieves the Grafana from the indexer for a given namespace and name.
func (s grafanaNamespaceLister) Get(name string) (*v1.Grafana, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("grafana"), name)
	}
	return obj.(*v1.Grafana), nil
}
