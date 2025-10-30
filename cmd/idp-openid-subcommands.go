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
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	json "github.com/openstor/colorjson"
	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/mc/pkg/probe"
	"github.com/urfave/cli/v3"
)

var idpOpenidAddCmd = cli.Command{
	Name:         "add",
	Usage:        "Create an OpenID IDP server configuration",
	Action:       mainIDPOpenIDAdd,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_NAME] [CFG_PARAMS...]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Create a default OpenID IDP configuration (CFG_NAME is omitted).
     {{.Prompt}} {{.HelpName}} play/ \
          client_id=minio-client-app \
          client_secret=minio-client-app-secret \
          config_url="http://localhost:5556/dex/.well-known/openid-configuration" \
          scopes="openid,groups" \
          redirect_uri="http://127.0.0.1:10000/oauth_callback" \
          role_policy="consoleAdmin"
  2. Create OpenID IDP configuration named "dex_test".
     {{.Prompt}} {{.HelpName}} play/ dex_test \
          client_id=minio-client-app \
          client_secret=minio-client-app-secret \
          config_url="http://localhost:5556/dex/.well-known/openid-configuration" \
          scopes="openid,groups" \
          redirect_uri="http://127.0.0.1:10000/oauth_callback" \
          role_policy="consoleAdmin"
`,
}

func mainIDPOpenIDAdd(ctx context.Context, cmd *cli.Command) error {
	return mainIDPOpenIDAddOrUpdate(ctx, cmd, false)
}

func mainIDPOpenIDAddOrUpdate(ctx context.Context, cmd *cli.Command, update bool) error {
	if cmd.Args().Len() < 2 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	args := cmd.Args()

	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	cfgName := madmin.Default
	input := args.Slice()[1:]
	if !strings.Contains(args.Get(1), "=") {
		cfgName = args.Get(1)
		input = args.Slice()[2:]
	}

	inputCfg := strings.Join(input, " ")

	restart, e := client.AddOrUpdateIDPConfig(globalContext, madmin.OpenidIDPCfg, cfgName, inputCfg, update)
	fatalIf(probe.NewError(e), "Unable to add OpenID IDP config to server")

	// Print set config result
	printMsg(configSetMessage{
		targetAlias: aliasedURL,
		restart:     restart,
	})

	return nil
}

var idpOpenidUpdateCmd = cli.Command{
	Name:         "update",
	Usage:        "Update an OpenID IDP configuration",
	Action:       mainIDPOpenIDUpdate,
	Before:       setGlobalsFromContext,
	OnUsageError: onUsageError,
	Flags:        globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_NAME] [CFG_PARAMS...]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Update the default OpenID IDP configuration (CFG_NAME is omitted).
     {{.Prompt}} {{.HelpName}} play/
          scopes="openid,groups" \
          role_policy="consoleAdmin"
  2. Update configuration for OpenID IDP configuration named "dex_test".
     {{.Prompt}} {{.HelpName}} play/ dex_test \
          scopes="openid,groups" \
          role_policy="consoleAdmin"
`,
}

func mainIDPOpenIDUpdate(ctx context.Context, cmd *cli.Command) error {
	return mainIDPOpenIDAddOrUpdate(ctx, cmd, true)
}

var idpOpenidRemoveCmd = cli.Command{
	Name:         "remove",
	Aliases:      []string{"rm"},
	Usage:        "remove OpenID IDP server configuration",
	Action:       mainIDPOpenIDRemove,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_NAME]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Remove the default OpenID IDP configuration (CFG_NAME is omitted).
     {{.Prompt}} {{.HelpName}} play/
  2. Remove OpenID IDP configuration named "dex_test".
     {{.Prompt}} {{.HelpName}} play/ dex_test
`,
}

func mainIDPOpenIDRemove(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 || cmd.Args().Len() > 2 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	args := cmd.Args()

	var cfgName string
	if args.Len() == 2 {
		cfgName = args.Get(1)
	}
	return idpRemove(ctx, cmd, true, cfgName)
}

func idpRemove(ctx context.Context, cmd *cli.Command, isOpenID bool, cfgName string) error {
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	idpType := madmin.LDAPIDPCfg
	if isOpenID {
		idpType = madmin.OpenidIDPCfg
	}

	restart, e := client.DeleteIDPConfig(globalContext, idpType, cfgName)
	fatalIf(probe.NewError(e), "Unable to remove %s IDP config '%s'", idpType, cfgName)

	printMsg(configSetMessage{
		targetAlias: aliasedURL,
		restart:     restart,
	})

	return nil
}

var idpOpenidListCmd = cli.Command{
	Name:         "list",
	Aliases:      []string{"ls"},
	Usage:        "list OpenID IDP server configuration(s)",
	Action:       mainIDPOpenIDList,
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
  1. List configurations for OpenID IDP.
     {{.Prompt}} {{.HelpName}} play/
`,
}

func mainIDPOpenIDList(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	return idpListCommon(ctx, cmd, true)
}

func idpListCommon(ctx context.Context, cmd *cli.Command, isOpenID bool) error {
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	idpType := madmin.LDAPIDPCfg
	if isOpenID {
		idpType = madmin.OpenidIDPCfg
	}
	result, e := client.ListIDPConfig(globalContext, idpType)
	fatalIf(probe.NewError(e), "Unable to list IDP config for '%s'", idpType)

	printMsg(idpCfgList(result))

	return nil
}

type idpCfgList []madmin.IDPListItem

func (i idpCfgList) JSON() string {
	bs, e := json.MarshalIndent(i, "", "  ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(bs)
}

func (i idpCfgList) String() string {
	maxNameWidth := len("Name")
	maxRoleARNWidth := len("RoleArn")
	for _, item := range i {
		name := item.Name
		if name == "_" {
			name = "(default)" // for the un-named config, don't show `_`
		}
		if maxNameWidth < len(name) {
			maxNameWidth = len(name)
		}
		if maxRoleARNWidth < len(item.RoleARN) {
			maxRoleARNWidth = len(item.RoleARN)
		}
	}
	enabledWidth := 5
	// Add 2 for padding
	maxNameWidth += 2
	maxRoleARNWidth += 2

	enabledColStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		PaddingLeft(1).
		PaddingRight(1).
		Width(enabledWidth)
	nameColStyle := lipgloss.NewStyle().
		Align(lipgloss.Right).
		PaddingLeft(1).
		PaddingRight(1).
		Width(maxNameWidth)
	arnColStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color("#04B575")). // green
		Width(maxRoleARNWidth)

	styles := []lipgloss.Style{enabledColStyle, nameColStyle, arnColStyle}

	headers := []string{"On?", "Name", "RoleARN"}
	headerRow := []string{}

	// Override some style settings for the header
	for ii, hdr := range headers {
		styl := styles[ii]
		headerRow = append(headerRow,
			styl.Bold(true).
				Foreground(lipgloss.Color("#6495ed")). // green
				Align(lipgloss.Center).
				Render(hdr),
		)
	}

	lines := []string{strings.Join(headerRow, "")}

	enabledOff := "🔴"
	enabledOn := "🟢"

	for _, item := range i {
		enabled := enabledOff
		if item.Enabled {
			enabled = enabledOn
		}

		line := []string{
			styles[0].Render(enabled),
			styles[1].Render(item.Name),
			styles[2].Render(item.RoleARN),
		}
		if item.Name == "_" {
			// For default config, don't display `_` and make it look faint.
			styl := styles[1]
			line[1] = styl.Faint(true).
				Render("(default)")
		}
		lines = append(lines, strings.Join(line, ""))
	}

	boxContent := strings.Join(lines, "\n")
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder())

	return boxStyle.Render(boxContent)
}

var idpOpenidInfoCmd = cli.Command{
	Name:         "info",
	Usage:        "get OpenID IDP server configuration info",
	Action:       mainIDPOpenIDInfo,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_NAME]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Get configuration info on the default OpenID IDP configuration (CFG_NAME is omitted).
     {{.Prompt}} {{.HelpName}} play/
  2. Get configuration info on OpenID IDP configuration named "dex_test".
     {{.Prompt}} {{.HelpName}} play/ dex_test
`,
}

func mainIDPOpenIDInfo(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 || cmd.Args().Len() > 2 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	args := cmd.Args()
	var cfgName string
	if args.Len() == 2 {
		cfgName = args.Get(1)
	}

	return idpInfo(ctx, cmd, true, cfgName)
}

func idpInfo(ctx context.Context, cmd *cli.Command, isOpenID bool, cfgName string) error {
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	idpType := madmin.LDAPIDPCfg
	if isOpenID {
		idpType = madmin.OpenidIDPCfg
	}

	result, e := client.GetIDPConfig(globalContext, idpType, cfgName)
	fatalIf(probe.NewError(e), "Unable to get %s IDP config from server", idpType)

	// Print set config result
	printMsg(idpConfig(result))

	return nil
}

type idpConfig madmin.IDPConfig

func (i idpConfig) JSON() string {
	bs, e := json.MarshalIndent(i, "", "  ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(bs)
}

func (i idpConfig) String() string {
	if len(i.Info) == 0 {
		return "Not configured."
	}

	enableStr := "on"
	// Determine required width for key column.
	fieldColWidth := 0
	for _, kv := range i.Info {
		if fieldColWidth < len(kv.Key) {
			fieldColWidth = len(kv.Key)
		}
		if kv.Key == "enable" {
			enableStr = kv.Value
		}
	}
	// Add 1 for the colon-suffix in each entry.
	fieldColWidth++

	fieldColStyle := lipgloss.NewStyle().
		Width(fieldColWidth).
		Foreground(lipgloss.Color("#04B575")). // green
		Bold(true).
		Align(lipgloss.Right)
	valueColStyle := lipgloss.NewStyle().
		PaddingLeft(1).
		Align(lipgloss.Left)
	envMarkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("201")). // pinkish-red
		PaddingLeft(1)

	var lines []string
	lines = append(lines, fmt.Sprintf("%s%s",
		fieldColStyle.Render("enable:"),
		valueColStyle.Render(enableStr),
	))
	for _, kv := range i.Info {
		if kv.Key == "enable" {
			continue
		}
		envStr := ""
		if kv.IsCfg && kv.IsEnv {
			envStr = " (environment)"
		}
		lines = append(lines, fmt.Sprintf("%s%s%s",
			fieldColStyle.Render(kv.Key+":"),
			valueColStyle.Render(kv.Value),
			envMarkStyle.Render(envStr),
		))
	}

	boxContent := strings.Join(lines, "\n")

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder())

	return boxStyle.Render(boxContent)
}

var idpOpenidEnableCmd = cli.Command{
	Name:         "enable",
	Usage:        "enable an OpenID IDP server configuration",
	Action:       mainIDPOpenIDEnable,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_NAME]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Enable the default OpenID IDP configuration (CFG_NAME is omitted).
     {{.Prompt}} {{.HelpName}} play/
  2. Enable OpenID IDP configuration named "dex_test".
     {{.Prompt}} {{.HelpName}} play/ dex_test
`,
}

func mainIDPOpenIDEnable(ctx context.Context, cmd *cli.Command) error {
	isOpenID, enable := true, true
	return idpEnableDisable(ctx, cmd, isOpenID, enable)
}

func mainIDPOpenIDDisable(ctx context.Context, cmd *cli.Command) error {
	isOpenID, enable := true, false
	return idpEnableDisable(ctx, cmd, isOpenID, enable)
}

func idpEnableDisable(ctx context.Context, cmd *cli.Command, isOpenID, enable bool) error {
	if cmd.Args().Len() < 1 || cmd.Args().Len() > 2 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}

	args := cmd.Args()
	cfgName := madmin.Default
	if args.Len() == 2 {
		cfgName = args.Get(1)
	}
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	idpType := madmin.LDAPIDPCfg
	if isOpenID {
		idpType = madmin.OpenidIDPCfg
	}

	configBody := "enable="
	if !enable {
		configBody = "enable=off"
	}

	restart, e := client.AddOrUpdateIDPConfig(globalContext, idpType, cfgName, configBody, true)
	fatalIf(probe.NewError(e), "Unable to remove %s IDP config '%s'", idpType, cfgName)

	printMsg(configSetMessage{
		targetAlias: aliasedURL,
		restart:     restart,
	})

	return nil
}

var idpOpenidDisableCmd = cli.Command{
	Name:         "disable",
	Usage:        "Disable an OpenID IDP server configuration",
	Action:       mainIDPOpenIDDisable,
	Before:       setGlobalsFromContext,
	Flags:        globalFlags,
	OnUsageError: onUsageError,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [CFG_NAME]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Disable the default OpenID IDP configuration (CFG_NAME is omitted).
     {{.Prompt}} {{.HelpName}} play/
  2. Disable OpenID IDP configuration named "dex_test".
     {{.Prompt}} {{.HelpName}} play/ dex_test
`,
}
