package config

const (
	OsConfigFile    = "/etc/k3os-config.yml"
	CloudConfigDir  = "/var/lib/k3os/conf/cloud-config.d"
	CloudConfigFile = "/var/lib/k3os/conf/cloud-config.yml"
)

var (
	OSVersion   string
	OSBuildDate string
	Additional  = []string{
		"EXTRA_CMDLINE",
	}
)

type CloudConfig struct {
	Hostname string     `yaml:"hostname,omitempty"`
	K3OS     K3OSConfig `yaml:"k3os,omitempty"`
}

type Defaults struct {
	Hostname string `yaml:"hostname,omitempty"`
}

type K3OSConfig struct {
	Defaults Defaults          `yaml:"defaults,omitempty"`
	Modules  []string          `yaml:"modules,omitempty"`
	SSH      SSHConfig         `yaml:"ssh,omitempty"`
	Sysctl   map[string]string `yaml:"sysctl,omitempty"`
	Upgrade  UpgradeConfig     `yaml:"upgrade,omitempty"`
}

type SSHConfig struct {
	Address        string   `yaml:"address,omitempty"`
	AuthorizedKeys []string `yaml:"authorized_keys,omitempty"`
	Daemon         bool     `yaml:"daemon,omitempty"`
	Port           int      `yaml:"port,omitempty"`
}

type UpgradeConfig struct {
	URL      string `yaml:"url,omitempty"`
	Image    string `yaml:"image,omitempty"`
	Rollback string `yaml:"rollback,omitempty"`
	Policy   string `yaml:"policy,omitempty"`
}
