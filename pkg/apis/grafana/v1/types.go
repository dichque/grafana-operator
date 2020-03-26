package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Grafana describes a Grafana resource
type Grafana struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec GrafanaSpec `json:"spec"`

	// Status is custom resource status
	// Status GrafanaStatus `json:"status,omitempty"`
}

// GrafanaSpec is the spec for a Grafana resource
type GrafanaSpec struct {
	// Message and SomeValue are example custom spec fields
	//
	// this is where you would put your custom resource data
	Image    string `json:"image"`
	Replicas *int32 `json:"replicas"`
}

// // GrafanaStatus defines the observed state of At
// type GrafanaStatus struct {
// 	// Message represents status of the customr resource
// 	Message string `json:"message,omitempty"`
// }

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaList is a list of Grafana resources
type GrafanaList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []Grafana `json:"items"`
}
