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
	"strings"

	"github.com/fatih/color"
	json "github.com/openstor/colorjson"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var adminConfigEnvFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "env",
		Usage: "list all the env only help",
	},
}

var adminConfigResetCmd = cli.Command{
	Name:         "reset",
	Usage:        "interactively reset a config key parameters",
	Before:       setGlobalsFromContext,
	Action:       mainAdminConfigReset,
	OnUsageError: onUsageError,
	Flags:        append(adminConfigEnvFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CONFIG-KEY...]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Reset MQTT notifcation target 'name1' settings to default values.
     {{.Prompt}} {{.HelpName}} myminio/ notify_mqtt:name1
  2. Reset compression's 'extensions' setting to default value.
     {{.Prompt}} {{.HelpName}} myminio/ compression extensions
  3. Reset site name and site region to default values.
     {{.Prompt}} {{.HelpName}} myminio/ site name region
`,
}

// configResetMessage container to hold locks information.
type configResetMessage struct {
	Status      string `json:"status"`
	targetAlias string
	key         string
	restart     bool
}

// String colorized service status message.
func (u configResetMessage) String() (msg string) {
	msg += console.Colorize("ResetConfigSuccess",
		fmt.Sprintf("'%s' is successfully reset.", u.key))
	if u.restart {
		suggestion := fmt.Sprintf("mc admin service restart %s", u.targetAlias)
		msg += console.Colorize("ResetConfigSuccess",
			fmt.Sprintf("\nPlease restart your server with `%s`.", suggestion))
	}
	return
}

// JSON jsonified service status message.
func (u configResetMessage) JSON() string {
	u.Status = "success"
	statusJSONBytes, e := json.MarshalIndent(u, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(statusJSONBytes)
}

// checkAdminConfigResetSyntax - validate all the passed arguments
func checkAdminConfigResetSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if !args.Present() {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

// main config set function
func mainAdminConfigReset(ctx context.Context, cmd *cli.Command) error {
	// Check command arguments
	checkAdminConfigResetSyntax(ctx, cmd)

	// Reset color preference of command outputs
	console.SetColor("ResetConfigSuccess", color.New(color.FgGreen, color.Bold))
	console.SetColor("ResetConfigFailure", color.New(color.FgRed, color.Bold))

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	if args.Len() == 1 {
		// Call get config API
		hr, e := client.HelpConfigKV(ctx, "", "", cmd.Bool("env"))
		fatalIf(probe.NewError(e), "Unable to get help for the sub-system")

		// Print
		printMsg(configHelpMessage{
			Value:   hr,
			envOnly: cmd.Bool("env"),
		})

		return nil
	}

	input := strings.Join(args.Tail(), " ")
	// Check if user has attempted to set values
	for _, k := range args.Tail() {
		if strings.Contains(k, "=") {
			e := fmt.Errorf("new settings may not be provided for sub-system keys")
			fatalIf(probe.NewError(e), "Unable to reset '%s' on the server", args.Tail()[0])
		}
	}

	// Call reset config API
	restart, e := client.DelConfigKV(ctx, input)
	fatalIf(probe.NewError(e), "Unable to reset '%s' on the server", input)

	// Print set config result
	printMsg(configResetMessage{
		targetAlias: aliasedURL,
		restart:     restart,
		key:         input,
	})

	return nil
}
