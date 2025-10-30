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

	"github.com/openstor/mc/pkg/probe"
	"github.com/urfave/cli/v3"
)

var adminTierVerifyCmd = cli.Command{
	Name:         "verify",
	Usage:        "verifies if remote tier configuration is valid",
	Action:       mainAdminTierVerify,
	Hidden:       true,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET NAME

NAME:
  Name of remote tier target. e.g WARM-TIER

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Verify if a tier config is valid.
     {{.Prompt}} {{.HelpName}} myminio WARM-TIER
`,
}

func mainAdminTierVerify(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	nArgs := args.Len()
	if nArgs < 2 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}
	if nArgs != 2 {
		fatalIf(errInvalidArgument().Trace(args.Slice()...),
			"Incorrect number of arguments for tier verify command.")
	}

	aliasedURL := args.Get(0)
	tierName := args.Get(1)
	if tierName == "" {
		fatalIf(errInvalidArgument(), "Tier name can't be empty")
	}

	// Create a new MinIO Admin Client
	client, cerr := newAdminClient(aliasedURL)
	fatalIf(cerr, "Unable to initialize admin connection.")

	e := client.VerifyTier(globalContext, tierName)
	fatalIf(probe.NewError(e).Trace(args.Slice()...), "Unable to verify remote tier target")

	printMsg(&tierMessage{
		op:       "verify",
		Status:   "success",
		TierName: tierName,
	})
	return nil
}
