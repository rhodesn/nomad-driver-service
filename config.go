package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/nomad/plugins/shared/hclspec"
)

var (
	// configSpec is the specification of the plugin's configuration
	// this is used to validate the configuration specified for the plugin
	// on the client.
	// this is not global, but can be specified on a per-client basis.
	configSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		// TODO: define plugin's agent configuration schema.
		//
		// The schema should be defined using HCL specs and it will be used to
		// validate the agent configuration provided by the user in the
		// `plugin` stanza (https://www.nomadproject.io/docs/configuration/plugin.html).
		//
		// For example, for the schema below a valid configuration would be:
		//
		//   plugin "hello-driver-plugin" {
		//     config {
		//       shell = "fish"
		//     }
		//   }
		"entrypoint": hclspec.NewDefault(
			hclspec.NewAttr("entrypoint", "string", false),
			hclspec.NewLiteral(`"systemd-run"`),
		),
	})

	// taskConfigSpec is the specification of the plugin's configuration for
	// a task
	// this is used to validated the configuration specified for the plugin
	// when a job is submitted.
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		// TODO: define plugin's task configuration schema
		"command":     hclspec.NewAttr("command", "string", true),
		"args":        hclspec.NewAttr("args", "list(string)", false),
		"description": hclspec.NewAttr("description", "string", false),
		"environment": hclspec.NewAttr("environment", "list(map(string))", false),
		"properties":  hclspec.NewAttr("properties", "list(map(string))", false),
		"unitname":    hclspec.NewAttr("unitname", "string", true),
		"slice":       hclspec.NewAttr("slice", "string", false),
		"dynamic_user": hclspec.NewDefault(
			hclspec.NewAttr("dynamic_user", "bool", false),
			hclspec.NewLiteral("true"),
		),
		"user": hclspec.NewAttr("user", "string", false),
	})
)

// Config contains configuration information for the plugin
type PluginConfig struct {
	// TODO: create decoded plugin configuration struct
	Entrypoint string `codec:"entrypoint"`
}

func (pc *PluginConfig) validate() error {
	if pc.Entrypoint != "systemd-run" && pc.Entrypoint != "systemctl" {
		return fmt.Errorf("couldn't find entrypoint %s", pc.Entrypoint)
	}

	return nil
}

// TaskConfig contains configuration information for a task that runs with
// this plugin
type TaskConfig struct {
	// Command
	Command string `codec:"command"`

	// Args
	Args []string `codec:"args"`

	// Description
	Description string `codec:"description"`

	// Environment
	Environment map[string]string `codec:"environment"`

	// Properties
	Properties map[string]string `codec:"properties"`

	// UnitName
	UnitName string `codec:"unitname"`

	// Slice
	Slice string `codec:"slice"`

	// DynamicUser
	DynamicUser bool `codec:"dynamic_user"`

	// User
	User string `codec:"user"`
}

func (tc *TaskConfig) validate() error {
	if tc.UnitName == "" {
		return errors.New("unitname config is empty")
	}

	if tc.DynamicUser && tc.User != "" {
		return errors.New("dynamic_user and user are mutually exclusive")
	}

	return nil
}
