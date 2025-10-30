// Copyright (c) 2015-2025 MinIO, Inc.
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

var idpLdapAccesskeySTSRevokeCmd = cli.Command{
	Name:         "sts-revoke",
	Usage:        "revokes all STS accounts or specified types for the specified user",
	Action:       mainIdpLdapAccesskeySTSRevoke,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        append(adminAccesskeySTSRevokeFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} ALIAS USER [--all | --token-type TOKEN_TYPE]

  Exactly one of --all or --token-type must be specified.

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Revoke all STS accounts for LDAP user 'bobfisher'
	 {{.Prompt}} {{.HelpName}} myminio uid=bobfisher,ou=people,ou=hwengg,dc=min,dc=io --all

  2. Revoke all STS accounts for LDAP user 'bobfisher' (alt)
	 {{.Prompt}} {{.HelpName}} myminio bobfisher --all

  3. Revoke STS accounts of a token type 'app-1' for user 'user1'
	 {{.Prompt}} {{.HelpName}} myminio user1 --token-type app-1

  4. Revoke all STS accounts for the authenticated user (must be LDAP service account)
	 {{.Prompt}} {{.HelpName}} myminio --self

  5. Revoke STS accounts of a token type 'app-1' for the authenticated user (must be LDAP service account)
	 {{.Prompt}} {{.HelpName}} myminio --self --token-type app-1
`,
}

// mainIdpLdapAccesskeySTSRevoke is the handle for "mc idp ldap accesskey sts-revoke" command.
func mainIdpLdapAccesskeySTSRevoke(ctx context.Context, cmd *cli.Command) error {
	checkSTSRevokeSyntax(ctx, cmd)

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)
	user := args.Get(1) // will be empty if --self flag is set
	tokenRevokeType := cmd.String("token-type")
	fullRevoke := cmd.Bool("all")

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	e := client.RevokeTokens(globalContext, madmin.RevokeTokensReq{
		User:            user,
		TokenRevokeType: tokenRevokeType,
		FullRevoke:      fullRevoke,
	})
	fatalIf(probe.NewError(e).Trace(args.Slice()...), "Unable to revoke tokens for %s", user)

	printMsg(stsRevokeMessage{
		User:            user,
		TokenRevokeType: tokenRevokeType,
	})

	return nil
}
