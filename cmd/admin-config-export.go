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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	json "github.com/openstor/colorjson"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var adminConfigExportCmd = cli.Command{
	Name:         "export",
	Usage:        "export all config keys to STDOUT",
	Before:       setGlobalsFromContext,
	Action:       mainAdminConfigExport,
	OnUsageError: onUsageError,
	Flags:        globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  The output includes environment variables set on the server. These cannot be overridden from the client.

  1. Export the current config from MinIO server
     {{.Prompt}} {{.HelpName}} play/ > config.txt
`,
}

// configExportMessage container to hold locks information.
type configExportMessage struct {
	Status string `json:"status"`
	Value  []byte `json:"value"`
}

// String colorized service status message.
func (u configExportMessage) String() string {
	console.SetColor("EnvVar", color.New(color.FgYellow))
	bio := bufio.NewReader(bytes.NewReader(u.Value))
	var lines []string
	for {
		s, e := bio.ReadString('\n')
		// Make lines displaying environment variables bold.
		if strings.HasPrefix(s, "# MINIO_") {
			s = strings.TrimPrefix(s, "# ")
			parts := strings.SplitN(s, "=", 2)
			s = fmt.Sprintf("# %s=%s", console.Colorize("EnvVar", parts[0]), parts[1])
			lines = append(lines, s)
		} else {
			lines = append(lines, s)
		}
		if e == io.EOF {
			break
		}
		fatalIf(probe.NewError(e), "Unable to marshal to string.")
	}
	return strings.Join(lines, "")
}

// JSON jsonified service status Message message.
func (u configExportMessage) JSON() string {
	u.Status = "success"
	statusJSONBytes, e := json.MarshalIndent(u, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(statusJSONBytes)
}

// checkAdminConfigExportSyntax - validate all the passed arguments
func checkAdminConfigExportSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if !args.Present() || args.Len() > 1 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

func mainAdminConfigExport(ctx context.Context, cmd *cli.Command) error {
	checkAdminConfigExportSyntax(ctx, cmd)

	// Export the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)

	// Create a new MinIO Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Unable to initialize admin connection.")

	// Call get config API
	buf, e := client.GetConfig(ctx)
	fatalIf(probe.NewError(e), "Unable to get server config")

	// Print
	printMsg(configExportMessage{
		Value: buf,
	})

	return nil
}
