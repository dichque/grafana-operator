package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/dichque/grafana-operator/pkg/apis/grafana"
)

// SchemeGroupVersion is the identifier for the API which includes
// the name of the group and the version of the API
var SchemeGroupVersion = schema.GroupVersion{
	Group:   Grafana.GroupName,
	Version: "v1",
}

// AddToScheme create a SchemeBuilder which uses functions to add types to
// the scheme
var AddToScheme = runtime.NewSchemeBuilder(addKnownTypes).AddToScheme

// Resource exported
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds our types to the API scheme by registering
// Grafana and GrafanaList
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&Grafana{},
		&GrafanaList{},
	)

	// register the type in the scheme
	meta_v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
