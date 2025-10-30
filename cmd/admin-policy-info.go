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
	"os"

	"github.com/fatih/color"
	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var policyInfoFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "policy-file, f",
		Usage: "additionally (over-)write policy JSON to given file",
	},
}

var adminPolicyInfoCmd = &cli.Command{
	Name:         "info",
	Usage:        "show info on an IAM policy",
	Action:       mainAdminPolicyInfo,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        append(policyInfoFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET POLICYNAME [OPTIONS...]

POLICYNAME:
  Name of the policy on the MinIO server.

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Show information on a given policy.
     {{.Prompt}} {{.HelpName}} myminio writeonly

  2. Show information on a given policy and write the policy JSON content to /tmp/policy.json.
     {{.Prompt}} {{.HelpName}} myminio writeonly --policy-file /tmp/policy.json
`,
}

// checkAdminPolicyInfoSyntax - validate all the passed arguments
func checkAdminPolicyInfoSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if args.Len() != 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

func getPolicyInfo(client *madmin.AdminClient, policyName string) (*madmin.PolicyInfo, error) {
	pinfo, e := client.InfoCannedPolicy(globalContext, policyName)
	if e != nil {
		return nil, e
	}
	return pinfo, nil
}

// mainAdminPolicyInfo is the handler for "mc admin policy info" command.
func mainAdminPolicyInfo(ctx context.Context, cmd *cli.Command) error {
	checkAdminPolicyInfoSyntax(ctx, cmd)

	console.SetColor("PolicyMessage", color.New(color.FgGreen))
	console.SetColor("Policy", color.New(color.FgBlue))

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)
	policyName := args.Get(1)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection")

	pinfo, e := getPolicyInfo(client, policyName)
	fatalIf(probe.NewError(e).Trace(args.Slice()...), "Unable to fetch policy")

	policyFile := cmd.String("policy-file")
	if policyFile != "" {
		f, e := os.Create(policyFile)
		fatalIf(probe.NewError(e).Trace(args.Slice()...), "Could not open given policy file")

		_, e = f.Write(pinfo.Policy)
		fatalIf(probe.NewError(e).Trace(args.Slice()...), "Could not write to given policy file")
	}

	printMsg(userPolicyMessage{
		op:         "info",
		Policy:     policyName,
		PolicyInfo: *pinfo,
	})

	return nil
}
