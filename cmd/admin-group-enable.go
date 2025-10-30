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
	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var adminGroupEnableCmd = cli.Command{
	Name:         "enable",
	Usage:        "enable a group",
	Action:       mainAdminGroupEnableDisable,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET GROUPNAME

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Enable group 'allcents'.
     {{.Prompt}} {{.HelpName}} myminio allcents
`,
}

// checkAdminGroupEnableSyntax - validate all the passed arguments
func checkAdminGroupEnableSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if args.Len() != 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

// mainAdminGroupEnableDisable is the handle for "mc admin group enable|disable" command.
func mainAdminGroupEnableDisable(ctx context.Context, cmd *cli.Command) error {
	checkAdminGroupEnableSyntax(ctx, cmd)

	console.SetColor("GroupMessage", color.New(color.FgGreen))

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	group := args.Get(1)
	var status madmin.GroupStatus
	switch cmd.Name {
	case "enable":
		status = madmin.GroupEnabled
	case "disable":
		status = madmin.GroupDisabled
	default:
		fatalIf(errInvalidArgument().Trace(cmd.Name), "Invalid group status name")
	}
	e := client.SetGroupStatus(ctx, group, status)
	// Convert cli.Args to []string for Trace
	argsSlice := make([]string, args.Len())
	for i := 0; i < args.Len(); i++ {
		argsSlice[i] = args.Get(i)
	}
	fatalIf(probe.NewError(e).Trace(argsSlice...), "Unable set group status")

	printMsg(groupMessage{
		op:          cmd.Name,
		GroupName:   group,
		GroupStatus: string(status),
	})

	return nil
}
