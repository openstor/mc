// Copyright (c) 2015-2022 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	json "github.com/openstor/colorjson"
	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/openstor-go/v7/pkg/set"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

const (
	supportSuccessMsgTag = "SupportSuccessMessage"
	supportErrorMsgTag   = "SupportErrorMessage"
)

var supportGlobalFlags = append(globalFlags,
	&cli.BoolFlag{
		Name:   "dev",
		Usage:  "Development mode",
		Hidden: true,
	},
	&cli.BoolFlag{
		Name:  "airgap",
		Usage: "use in environments without network access to SUBNET (e.g. airgapped, firewalled, etc.)",
	},
)

var supportSubcommands = []*cli.Command{
	&supportRegisterCmd,
	&supportCallhomeCmd,
	&supportDiagCmd,
	&supportPerfCmd,
	&supportInspectCmd,
	&supportProfileCmd,
	&supportTopCmd,
	&supportProxyCmd,
	&supportUploadCmd,
}

var supportCmd = cli.Command{
	Name:            "support",
	Usage:           "support related commands",
	Action:          mainSupport,
	Before:          setGlobalsFromContext,
	Flags:           globalFlags,
	Commands:        supportSubcommands,
	HideHelpCommand: true,
}

func toggleCmdArgs() set.StringSet {
	return set.CreateStringSet("enable", "disable", "status")
}

func validateToggleCmdArg(arg string) error {
	valid := toggleCmdArgs()
	if !valid.Contains(arg) {
		return fmt.Errorf("Invalid argument '%s'. Must be one of %v", arg, valid)
	}
	return nil
}

func checkToggleCmdSyntax(ctx context.Context, cmd *cli.Command) (string, string) {
	if cmd.Args().Len() != 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}

	arg := cmd.Args().Get(0)
	aliasedURL := cmd.Args().Get(1)
	fatalIf(probe.NewError(validateToggleCmdArg(arg)), "Invalid arguments.")

	alias, _ := url2Alias(aliasedURL)

	return alias, arg
}

func setSuccessMessageColor() {
	console.SetColor(supportSuccessMsgTag, color.New(color.FgGreen, color.Bold))
}

func setErrorMessageColor() {
	console.SetColor(supportErrorMsgTag, color.New(color.FgYellow, color.Italic))
}

func featureStatusStr(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func validateClusterRegistered(alias string, cmdTalksToSubnet bool) string {
	// Non-registered execution allowed only in following scenarios
	// command doesn't talk to subnet: dev mode (`--dev` passed)
	// command talks to subnet: dev+airgapped mode (both `--dev` and `--airgap` passed)
	requireRegistration := !GlobalDevMode
	if cmdTalksToSubnet {
		requireRegistration = !GlobalDevMode || !globalAirgapped
	}

	apiKey, e := getSubnetAPIKey(alias)
	if requireRegistration {
		fatalIf(probe.NewError(e), "")
	}

	return apiKey
}

// isFeatureEnabled - checks if a feature is enabled in MinIO config
// To be used with configs that can be switched on/off using the `enable` key
// e.g. subSys = logger_webhook, target = logger_webhook:subnet
// Returns true if any of the following is true
// - `enable` is set to `on`
// - `enable` key is not found
// Returns false if any of the following is true
// - given subsystem is not supported by the version of MinIO
// - the given target doesn't exist in the config
// - `enable` is set to `off`
func isFeatureEnabled(alias, subSys, target string) bool {
	client, err := newAdminClient(alias)
	// Create a new MinIO Admin Client
	fatalIf(err, "Unable to initialize admin connection.")

	if !minioConfigSupportsSubSys(client, subSys) {
		return false
	}

	scfgs, e := getMinIOSubSysConfig(client, subSys)
	fatalIf(probe.NewError(e), "Unable to list access keys.")
	if e != nil {
		// Ignore error if the given target doesn't exist
		// e.g. logger_webhook:subnet doesn't exist when
		// pushing logs to SUBNET has not been enabled
		if e.Error() == fmt.Sprintf("sub-system target '%s' doesn't exist", target) {
			return false
		}

		fatalIf(probe.NewError(e), fmt.Sprintf("Unable to get server config for '%s'", subSys))
	}

	if target == madmin.Default {
		target = ""
	}
	for _, scfg := range scfgs {
		if scfg.Target == target {
			enable, found := scfg.Lookup(madmin.EnableKey)
			if !found {
				// if `enable` key is not found, it means that `enable=on`
				return true
			}
			return enable == madmin.EnableOn
		}
	}
	return false
}

func toJSON(obj interface{}) string {
	jsonBytes, e := json.MarshalIndent(obj, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(jsonBytes)
}

// mainSupport is the handle for "mc support" command.
func mainSupport(ctx context.Context, cmd *cli.Command) error {
	// Convert []*cli.Command to []cli.Command for compatibility
	var subCmds []cli.Command
	for _, c := range supportSubcommands {
		subCmds = append(subCmds, *c)
	}
	commandNotFound(ctx, cmd, subCmds)
	return nil
	// Sub-commands like "register", "callhome", "diagnostics" have their own main.
}
