package config

import (
	"fmt"
	"os"
	"strconv"
)

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
	K3sArgs        []string          `json:"k3sArgs,omitempt"`
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
	Initcmd           []string   `json:"initCmd,omitempty"`
}

type File struct {
	Encoding           string `json:"encoding"`
	Content            string `json:"content"`
	Owner              string `json:"owner"`
	Path               string `json:"path"`
	RawFilePermissions string `json:"permissions"`
}

func (f *File) Permissions() (os.FileMode, error) {
	if f.RawFilePermissions == "" {
		return os.FileMode(0644), nil
	}
	// parse string representation of file mode as integer
	perm, err := strconv.ParseInt(f.RawFilePermissions, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("unable to parse file permissions %q as integer", f.RawFilePermissions)
	}
	return os.FileMode(perm), nil
}
