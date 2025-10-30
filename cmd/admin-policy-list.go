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

var adminPolicyListCmd = &cli.Command{
	Name:         "list",
	Aliases:      []string{"ls"},
	Usage:        "list all IAM policies",
	Action:       mainAdminPolicyList,
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
  1. List all policies on MinIO server.
     {{.Prompt}} {{.HelpName}} myminio
`,
}

// checkAdminPolicyListSyntax - validate all the passed arguments
func checkAdminPolicyListSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if args.Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

// mainAdminPolicyList is the handle for "mc admin policy add" command.
func mainAdminPolicyList(ctx context.Context, cmd *cli.Command) error {
	checkAdminPolicyListSyntax(ctx, cmd)

	console.SetColor("PolicyName", color.New(color.FgBlue))

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	policies, e := client.ListCannedPolicies(globalContext)
	fatalIf(probe.NewError(e).Trace(args.Slice()...), "Unable to list policy")

	for k := range policies {
		printMsg(userPolicyMessage{
			op:     "list",
			Policy: k,
		})
	}
	return nil
}
