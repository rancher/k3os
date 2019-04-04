package network

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func CheckDhcpcd() error {
	cmd := exec.Command("/bin/sh", "-c", `ps aux | grep "dhcpcd"`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err
	}
	if len(bytesErr) != 0 {
		return errors.New(string(bytesErr))
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	if strings.Contains(string(bytes), "/sbin/dhcpcd") {
		return errors.New("dhcpcd process is still running")
	}
	return nil
}

func ReleaseDhcpcd(iface string) error {
	args := []string{"-k"}
	if iface != "" {
		args = append(args, iface)
	}
	cmd := exec.Command("/sbin/dhcpcd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func StartDhcpcd() error {
	cmd := exec.Command("/sbin/dhcpcd", "-f", "/etc/dhcpcd.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func StopDhcpcd() error {
	cmd := exec.Command("/sbin/dhcpcd", "-x")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
