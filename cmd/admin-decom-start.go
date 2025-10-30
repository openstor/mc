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
	"path/filepath"

	"github.com/fatih/color"
	json "github.com/openstor/colorjson"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var adminDecommissionStartCmd = cli.Command{
	Name:         "start",
	Usage:        "start decommissioning a pool",
	Action:       mainAdminDecommissionStart,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Start decommissioning a pool for removal.
     {{.Prompt}} {{.HelpName}} myminio/ http://server{5...8}/disk{1...4}
`,
}

// checkAdminDecommissionStartSyntax - validate all the passed arguments
func checkAdminDecommissionStartSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if args.Len() != 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

// startDecomMessage is container for make bucket success and failure messages.
type startDecomMessage struct {
	Status string `json:"status"`
	Pool   string `json:"pool"`
}

// String colorized construct a string message.
func (s startDecomMessage) String() string {
	return console.Colorize("DecomPool", "Decommission started successfully for `"+s.Pool+"`.")
}

// JSON jsonified decom message.
func (s startDecomMessage) JSON() string {
	startDecomBytes, e := json.MarshalIndent(s, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(startDecomBytes)
}

// mainAdminDecommissionStart is the handle for "mc admin decommission start" command.
func mainAdminDecommissionStart(ctx context.Context, cmd *cli.Command) error {
	checkAdminDecommissionStartSyntax(ctx, cmd)

	// Additional command speific theme customization.
	console.SetColor("DecomPool", color.New(color.FgGreen, color.Bold))

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)
	aliasedURL = filepath.Clean(aliasedURL)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	e := client.DecommissionPool(ctx, args.Get(1))
	// Convert cli.Args to []string for Trace
	argsSlice := make([]string, args.Len())
	for i := 0; i < args.Len(); i++ {
		argsSlice[i] = args.Get(i)
	}
	fatalIf(probe.NewError(e).Trace(argsSlice...), "Unable to start decommission on the specified pool")

	printMsg(startDecomMessage{
		Status: "success",
		Pool:   args.Get(1),
	})
	return nil
}
