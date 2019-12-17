package v1

// Copyright 2019 Rancher Labs, Inc.
// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/genericcondition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// ChannelLatestResolved indicates the value for latest has been polled and resolved
	ChannelLatestResolved = condition.Cond("LatestResolved")
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Channel represents a series of releases. The only protocol that a Channel
// must adhere to is GETing the URL will return a 302 with the last segment
// of the Location being the "latest" tag/release.
type Channel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChannelSpec   `json:"spec"`
	Status ChannelStatus `json:"status"`
}

// ChannelSpec represents the user-configurable details of a Channel
type ChannelSpec struct {
	URL string `json:"url,omitempty"` // e.g. https://github.com/rancher/k3os/releases/latest
}

// ChannelStatus represents the resulting state from processing Channel events.
type ChannelStatus struct {
	LatestVersion string                              `json:"latestVersion,omitempty"`
	Conditions    []genericcondition.GenericCondition `json:"conditions,omitempty"`
}

type DrainSpec struct {
	Timeout          *time.Duration `json:"timeout,omitempty"`
	GracePeriod      *int32         `json:"gracePeriod,omitempty"`
	DeleteLocalData  bool           `json:"deleteLocalData,omitempty"`
	IgnoreDaemonSets bool           `json:"ignoreDaemonSets,omitempty"`
	Force            bool           `json:"force,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpgradeSet represents a "JobSet" of NodeUpgrades
type UpgradeSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpgradeSetSpec   `json:"spec,omitempty"`
	Status UpgradeSetStatus `json:"status,omitempty"`
}

// UpgradeSetSpec represents the user-configurable details of an UpgradeSet
type UpgradeSetSpec struct {
	Concurrency uint64     `json:"concurrency,omitempty"`
	Version     string     `json:"version,omitempty"`
	Drain       *DrainSpec `json:"drain,omitempty"`
}

// UpgradeSetStatus represents the resulting state from processing UpgradeSet events.
type UpgradeSetStatus struct {
	Upgrades []string `json:"upgrades,omitempty"`
}

// NodeUpgradeStatus conditions
var (
	// NodeUpgradeScheduled indicates that the job has been scheduled.
	NodeUpgradeScheduled = condition.Cond("Scheduled")

	// NodeUpgradeCordoned indicates status of cordoning the node.
	NodeUpgradeCordoned = condition.Cond("Cordoned")

	// NodeUpgradeApplied indicates status of applying the upgrade.
	NodeUpgradeApplied = condition.Cond("Applied")

	// NodeUpgradeDrained indicates status of draining the node.
	NodeUpgradeDrained = condition.Cond("Drained")

	// NodeUpgradeRebooted indicates status of rebooting the node.
	NodeUpgradeRebooted = condition.Cond("Rebooted")
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeUpgrade represents a scheduled "Job" for performing an upgrade of a node.
type NodeUpgrade struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeUpgradeSpec   `json:"spec,omitempty"`
	Status NodeUpgradeStatus `json:"status,omitempty"`
}

// NodeUpgradeSpec represents the parameters necessary for creating the underlying Job.
type NodeUpgradeSpec struct {
	NodeName string     `json:"nodeName,omitempty"`
	Version  string     `json:"version,omitempty"`
	Drain    *DrainSpec `json:"drain,omitempty"`
}

// NodeUpgradeStatus represents the resulting state from processing NodeUpgrade events.
type NodeUpgradeStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`
}
