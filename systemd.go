package main

import (
	"os/exec"
	"regexp"

	systemd "github.com/coreos/go-systemd/v22/dbus"
	"github.com/godbus/dbus/v5"
)

var defaultProperties []systemd.Property = []systemd.Property{
	newSystemdProperty("NoNewPrivileges", true),
}

func newSystemdProperty(name string, value interface{}) systemd.Property {
	return systemd.Property{
		Name:  name,
		Value: dbus.MakeVariant(value),
	}
}

func systemdVersion() (string, error) {
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(systemctl, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`^systemd [0-9]{3,}`)
	return re.FindString(string(out)), nil
}
