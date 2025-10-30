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
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	json "github.com/openstor/colorjson"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var adminKMSKeyListCmd = &cli.Command{
	Name:         "list",
	Usage:        "request list of KMS master keys",
	Action:       mainAdminKMSKeyList,
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
  1. Get list of master keys from a MinIO server/cluster.
     $ {{.HelpName}} play
`,
}

// adminKMSKeyCmd is the handle for the "mc admin kms key" command.
func mainAdminKMSKeyList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	if args.Len() == 0 || args.Len() > 1 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}

	console.SetColor("KeyName", color.New(color.FgBlue))

	// Get the alias parameter from cli
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	keys, e := client.ListKeys(globalContext, "*")
	fatalIf(probe.NewError(e).Trace(args.Slice()...), "Unable to list KMS keys")

	var rows []table.Row
	kmsKeys := []string{}
	for idx, k := range keys {
		rows = append(rows, table.Row{idx + 1, k.Name})
		kmsKeys = append(kmsKeys, k.Name)
	}

	if globalJSON {
		printMsg(kmsKeysMsg{
			Status: "success",
			Target: aliasedURL,
			Keys:   kmsKeys,
		})
		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetColumnConfigs([]table.ColumnConfig{{Align: text.AlignCenter}})
	t.SetTitle("KMS Keys")
	t.AppendHeader(table.Row{"S N", "Name"})
	t.AppendRows(rows)
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}

type kmsKeysMsg struct {
	Status string   `json:"status"`
	Target string   `json:"target"`
	Keys   []string `json:"keys"`
}

func (k kmsKeysMsg) JSON() string {
	kmsBytes, e := json.MarshalIndent(k, "", "    ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(kmsBytes)
}

func (k kmsKeysMsg) String() string {
	return fmt.Sprintf("Keys: %s\n", k.Keys)
}
