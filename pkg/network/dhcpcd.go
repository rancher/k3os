package network

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func CheckDhcpcd() (bool, error) {
	cmd := exec.Command("/bin/sh", "-c", `ps aux | grep "dhcpcd"`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return false, err
	}
	if err := cmd.Start(); err != nil {
		return false, err
	}
	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return false, err
	}
	if len(bytesErr) != 0 {
		return false, errors.New(string(bytesErr))
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return false, err
	}
	if !strings.Contains(string(bytes), "/sbin/dhcpcd") {
		return false, nil
	}
	return true, nil
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

func StartDhcpcd(arguments []string) error {
	args := []string{"-f", "/etc/dhcpcd.conf"}
	if len(arguments) > 0 {
		args = append(args, arguments...)
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

func RequestDhcpcd(iface string) error {
	cmd := exec.Command("/sbin/dhcpcd", "-A4", iface)
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
