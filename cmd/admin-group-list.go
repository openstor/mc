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

	"github.com/fatih/color"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var adminGroupListCmd = cli.Command{
	Name:         "list",
	Aliases:      []string{"ls"},
	Usage:        "display list of groups",
	Action:       mainAdminGroupList,
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
  1. List all groups.
     {{.Prompt}} {{.HelpName}} myminio
`,
}

// checkAdminGroupListSyntax - validate all the passed arguments
func checkAdminGroupListSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if args.Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

// mainAdminGroupList is the handle for "mc admin group list" command.
func mainAdminGroupList(ctx context.Context, cmd *cli.Command) error {
	checkAdminGroupListSyntax(ctx, cmd)

	console.SetColor("GroupMessage", color.New(color.FgGreen))

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	gs, e := client.ListGroups(ctx)
	// Convert cli.Args to []string for Trace
	argsSlice := make([]string, args.Len())
	for i := 0; i < args.Len(); i++ {
		argsSlice[i] = args.Get(i)
	}
	fatalIf(probe.NewError(e).Trace(argsSlice...), "Unable to list groups")

	printMsg(groupMessage{
		op:     "list",
		Groups: gs,
	})

	return nil
}
