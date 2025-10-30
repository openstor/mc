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
	"os"

	"github.com/fatih/color"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var adminKMSCreateKeyCmd = &cli.Command{
	Name:         "create",
	Usage:        "creates a new master KMS key",
	Action:       mainAdminKMSCreateKey,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [KEY_NAME]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Create a new master key named 'my-key' default master key.
     $ {{.HelpName}} play my-key
`,
}

// adminKMSCreateKeyCmd is the handler for the "mc admin kms key create" command.
func mainAdminKMSCreateKey(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}

	client, err := newAdminClient(cmd.Args().Get(0))
	fatalIf(err, "Cannot get a configured admin connection.")

	keyID := cmd.Args().Get(1)
	e := client.CreateKey(globalContext, keyID)
	fatalIf(probe.NewError(e), "Failed to create master key")

	if term.IsTerminal(int(os.Stdout.Fd())) {
		console.Println(color.GreenString(fmt.Sprintf("Created master key `%s` successfully", keyID)))
	}
	return nil
}
