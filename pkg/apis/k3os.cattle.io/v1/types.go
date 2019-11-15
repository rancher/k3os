package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UpdateChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpdateChannelSpec   `json:"spec"`
	Status UpdateChannelStatus `json:"status"`
}

type UpdateChannelSpec struct {
	URL         string `json:"url,omitempty"`
	Version     string `json:"version,omitempty"`
	Concurrency int    `json:"concurrency,omitempty"`
}

type UpdateChannelStatus struct {
	Polling   string   `json:"polling,omitempty"`
	Upgrading []string `json:"upgrading,omitempty"`
}
