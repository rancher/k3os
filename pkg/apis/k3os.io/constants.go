package k3os

const (
	// LabelMode represents the mode of k3OS
	LabelMode = GroupName + `/mode`

	// LabelVersion represents the version of k3OS
	LabelVersion = GroupName + `/version`

	// LabelUpgradeChannel represents a k3OS upgrade channel (by name)
	LabelUpgradeChannel = `upgrade.` + GroupName + `/channel`

	// LabelUpgradeEnabled indicates that upgrades should be applied
	LabelUpgradeEnabled = `upgrade.` + GroupName + `/enabled`

	// LabelUpgradeNoDrain on a pod indicates that it should not be terminated when draining
	LabelUpgradeOperator = `upgrade.` + GroupName + `/operator`

	// LabelUpgradeVersion represents a k3OS upgrade target version
	LabelUpgradeVersion = `upgrade.` + GroupName + `/version`
)
