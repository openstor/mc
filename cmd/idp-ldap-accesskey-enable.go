// Copyright (c) 2015-2024 MinIO, Inc.
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

	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/mc/pkg/probe"
	"github.com/urfave/cli/v3"
)

var idpLdapAccesskeyEnableCmd = cli.Command{
	Name:         "enable",
	Usage:        "enable an access key",
	Action:       mainIDPLdapAccesskeyEnable,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] [TARGET]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Enable LDAP access key
	 {{.Prompt}} {{.HelpName}} myminio myaccesskey
`,
}

func mainIDPLdapAccesskeyEnable(ctx context.Context, cmd *cli.Command) error {
	return enableDisableAccesskey(ctx, cmd, true)
}

func enableDisableAccesskey(ctx context.Context, cmd *cli.Command, enable bool) error {
	if cmd.Args().Len() == 0 || cmd.Args().Len() > 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}

	args := cmd.Args()
	aliasedURL := args.Get(0)
	accessKey := args.Get(1)

	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	op := "disable"
	status := "off"
	if enable {
		op = "enable"
		status = "on"
	}

	e := client.UpdateServiceAccount(globalContext, accessKey, madmin.UpdateServiceAccountReq{
		NewStatus: status,
	})
	fatalIf(probe.NewError(e), "Unable to add service account.")

	m := accesskeyMessage{
		op:        op,
		Status:    "success",
		AccessKey: accessKey,
	}
	printMsg(m)

	return nil
}
