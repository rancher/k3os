package apply

import "github.com/rancher/k3os/pkg/config"

/*
type K3OS struct {
	Mode           string            `json:"mode,omitempty"`
	DataSources    []string          `json:"dataSources,omitempty"`
	Modules        []string          `json:"modules,omitempty"`
	Sysctls        map[string]string `json:"sysctls,omitempty"`
	DNSNameservers []string          `json:"dnsNameServers,omitempty"`
	DNSSearch      []string          `json:"dnsSearch,omitempty"`
	DNSOptions     []string          `json:"dnsOptions,omitempty"`
	Password       string            `json:"password,omitempty"`
	ServerURL      string            `json:"serverUrl,omitempty"`
	Token          string            `json:"token,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	Taints         []string          `json:"taints,omitempty"`
	Install        Install           `json:"install,omitempty"`
}

type Install struct {
	EFI       bool   `json:"efi,omitempty"`
	MSDOS     bool   `json:"msdos,omitempty"`
	Device    string `json:"device,omitempty"`
	ConfigURL string `json:"configUrl,omitempty"`
	Silent    bool   `json:"silent,omitempty"`
}

type CloudConfig struct {
	SSHAuthorizedKeys []string   `json:"sshAuthorizedKeys,omitempty"`
	WriteFiles        []File     `json:"writeFiles,omitempty"`
	Hostname          string     `json:"hostname,omitempty"`
	Mounts            [][]string `json:"mounts,omitempty"`
	K3OS              K3OS       `json:"k3os,omitempty"`
	Runcmd            []string   `json:"runCmd,omitempty"`
	Bootcmd           []string   `json:"bootCmd,omitempty"`
	Initcmd           []string   `json:"bootCmd,omitempty"`
}
*/

type applier func(cfg *config.CloudConfig) error

func runApplies(cfg *config.CloudConfig, appliers ...applier) error {

}

func RunApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplySSHKeysWithNet,
		ApplyWriteFiles,
		ApplyRuncmd,
	)
}

func ConfigApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyK3S,
	)
}

func BootApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyModules,
		ApplySysctls,
		ApplyHostname,
		//ApplyDNS,
		ApplyPassword,
		//ApplyMounts,
		ApplySSHKeys,
		ApplyK3S,
		ApplyWriteFiles,
		ApplyBootcmd,
	)
}

func InitApply(cfg *config.CloudConfig) error {
	return runApplies(cfg,
		ApplyModules,
		ApplySysctls,
		ApplyHostname,
		//ApplyDNS,
		ApplyWriteFiles,
		ApplyInitcmd,
	)
}
