package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

const (
	PhasePending = "PENDING"
	PhaseRunning = "RUNNING"
	PhaseDone    = "DONE"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Grafana describes a Grafana resource
type Grafana struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaSpec   `json:"spec"`
	Status GrafanaStatus `json:"status,omitempty"`
}

// GrafanaSpec is the spec for a Grafana resource
type GrafanaSpec struct {
	Image    string `json:"image"`
	Replicas *int32 `json:"replicas"`
}

// GrafanaStatus defines the observed state of At
type GrafanaStatus struct {
	Phase string `json:"phase,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaList is a list of Grafana resources
type GrafanaList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata,omitempty"`

	Items []Grafana `json:"items"`
}
