package v1

import (
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// GrafanaSpec is the spec for a grafana resource
type GrafanaSpec struct {
	Image         string `json:"image,omitempty"`
	Replicas      *int32 `json:"replicas,omitempty"`
	Username      string `json:"user,omitempty"`
	Password      string `json:"password,omitempty"`
	PrometheusURL string `json:"prometheus_url,omitempty"`
}

// GrafanaStatus defines the observed state of grafana custom resource
type GrafanaStatus struct {
	GStatus         v1.ConditionStatus `json:"gStatus,omitempty"`
	LastUpdatedTime meta_v1.Time       `json:"lastUpdatedTime,omitempty"`
	Conditions      []GrafanaCondition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaList is a list of Grafana resources
type GrafanaList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata,omitempty"`
	Items            []Grafana `json:"items"`
}

// ConditionType we track
type ConditionType string

const (
	// ConditionTypeGrafanaConfigMap tracks configmap
	ConditionTypeGrafanaConfigMap ConditionType = "GrafanaConfigMap"

	// ConditionTypeGrafanaDeployment tracks deployment
	ConditionTypeGrafanaDeployment ConditionType = "GrafanaDeployment"
)

// ConditionStatus we track
type ConditionStatus string

// These are valid condition status. "ConditionStatusTrue" means a resource is in the condition;
// "ConditionStatusFalse" means a resource is not in the condition; "ConditionStatusUnknown" means kubernetes
// can't decide if a resource is in the condition or not.
const (
	ConditionStatusTrue    ConditionStatus = "True"
	ConditionStatusFalse   ConditionStatus = "False"
	ConditionStatusUnknown ConditionStatus = "Unknown"
)

type ConditionReason string

const (
	ConditionReasonGrafanaConfigMapUpdate  ConditionReason = "ConfigMapUpdate"
	ConditionReasonGrafanaDeploymentUpdate ConditionReason = "DeploymentUpdate"

	ConditionReasonGrafanaConfigMapDelete  ConditionReason = "ConfigMapDelete"
	ConditionReasonGrafanaDeploymentDelete ConditionReason = "DeploymentDelete"
)

// GrafanaCondition defines the observed state of grafana custom resource
type GrafanaCondition struct {
	Type               ConditionType   `json:"type"`
	Status             ConditionStatus `json:"status"`
	Reason             ConditionReason `json:"reason"`
	Message            string          `json:"message,omitempty"`
	LastTransitionTime meta_v1.Time    `json:"lastTransitionTime,omitempty"`
}
