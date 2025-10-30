// Copyright (c) 2015-2023 MinIO, Inc.
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
	"errors"
	"strings"

	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/mc/pkg/probe"
	"github.com/urfave/cli/v3"
)

var idpLdapAddCmd = cli.Command{
	Name:         "add",
	Usage:        "Create an LDAP IDP server configuration",
	Action:       mainIDPLDAPAdd,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_PARAMS...]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Create LDAP IDentity Provider configuration.
     {{.Prompt}} {{.HelpName}} myminio/ \
          server_addr=myldapserver:636 \
          lookup_bind_dn=cn=admin,dc=min,dc=io \
          lookup_bind_password=somesecret \
          user_dn_search_base_dn=dc=min,dc=io \
          user_dn_search_filter="(uid=%s)" \
          group_search_base_dn=ou=swengg,dc=min,dc=io \
          group_search_filter="(&(objectclass=groupofnames)(member=%d))"
`,
}

func mainIDPLDAPAdd(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 2 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	args := cmd.Args()

	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	cfgName := madmin.Default
	input := args.Tail()
	if !strings.Contains(args.Get(1), "=") {
		cfgName = args.Get(1)
		input = args.Tail()[1:]
	}

	if cfgName != madmin.Default {
		fatalIf(probe.NewError(errors.New("all config parameters must be of the form \"key=value\"")),
			"Bad LDAP IDP configuration")
	}

	inputCfg := strings.Join(input, " ")

	restart, e := client.AddOrUpdateIDPConfig(globalContext, madmin.LDAPIDPCfg, cfgName, inputCfg, false)
	fatalIf(probe.NewError(e), "Unable to add LDAP IDP config to server")

	// Print set config result
	printMsg(configSetMessage{
		targetAlias: aliasedURL,
		restart:     restart,
	})

	return nil
}

var idpLdapUpdateCmd = cli.Command{
	Name:         "update",
	Usage:        "Update an LDAP IDP configuration",
	Action:       mainIDPLDAPUpdate,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_PARAMS...]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Update the LDAP IDP configuration.
     {{.Prompt}} {{.HelpName}} play/ \
          lookup_bind_dn=cn=admin,dc=min,dc=io \
          lookup_bind_password=somesecret
`,
}

func mainIDPLDAPUpdate(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 2 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	args := cmd.Args()

	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	cfgName := madmin.Default
	input := args.Tail()
	if !strings.Contains(args.Get(1), "=") {
		cfgName = args.Get(1)
		input = args.Tail()[1:]
	}

	if cfgName != madmin.Default {
		fatalIf(probe.NewError(errors.New("all config parameters must be of the form \"key=value\"")),
			"Bad LDAP IDP configuration")
	}

	inputCfg := strings.Join(input, " ")

	restart, e := client.AddOrUpdateIDPConfig(globalContext, madmin.LDAPIDPCfg, cfgName, inputCfg, true)
	fatalIf(probe.NewError(e), "Unable to update LDAP IDP configuration")

	// Print set config result
	printMsg(configSetMessage{
		targetAlias: aliasedURL,
		restart:     restart,
	})

	return nil
}

var idpLdapRemoveCmd = cli.Command{
	Name:         "remove",
	Aliases:      []string{"rm"},
	Usage:        "remove LDAP IDP server configuration",
	Action:       mainIDPLDAPRemove,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Remove the default LDAP IDP configuration.
     {{.Prompt}} {{.HelpName}} play/
`,
}

func mainIDPLDAPRemove(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	cfgName := madmin.Default
	return idpRemove(ctx, cmd, false, cfgName)
}

var idpLdapListCmd = cli.Command{
	Name:         "list",
	Aliases:      []string{"ls"},
	Usage:        "list LDAP IDP server configuration(s)",
	Action:       mainIDPLDAPList,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. List configurations for LDAP IDP.
     {{.Prompt}} {{.HelpName}} play/
`,
}

func mainIDPLDAPList(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	return idpListCommon(ctx, cmd, false)
}

var idpLdapInfoCmd = cli.Command{
	Name:         "info",
	Usage:        "get LDAP IDP server configuration info",
	Action:       mainIDPLDAPInfo,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Get configuration info on the LDAP IDP configuration.
     {{.Prompt}} {{.HelpName}} play/
`,
}

func mainIDPLDAPInfo(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	cfgName := madmin.Default
	return idpInfo(ctx, cmd, false, cfgName)
}

var idpLdapEnableCmd = cli.Command{
	Name:         "enable",
	Usage:        "manage LDAP IDP server configuration",
	Action:       mainIDPLDAPEnable,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Enable the LDAP IDP configuration.
     {{.Prompt}} {{.HelpName}} play/
`,
}

func mainIDPLDAPEnable(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	isOpenID, enable := false, true
	return idpEnableDisable(ctx, cmd, isOpenID, enable)
}

var idpLdapDisableCmd = cli.Command{
	Name:         "disable",
	Usage:        "Disable an LDAP IDP server configuration",
	Action:       mainIDPLDAPDisable,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Disable the default LDAP IDP configuration.
     {{.Prompt}} {{.HelpName}} play/
`,
}

func mainIDPLDAPDisable(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	isOpenID, enable := false, false
	return idpEnableDisable(ctx, cmd, isOpenID, enable)
}
