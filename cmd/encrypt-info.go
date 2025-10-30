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

	"github.com/fatih/color"
	json "github.com/openstor/colorjson"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var encryptInfoCmd = cli.Command{
	Name:         "info",
	Usage:        "show bucket encryption status",
	Action:       mainEncryptInfo,
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
  1. Display bucket encryption status for bucket "mybucket".
     {{.Prompt}} {{.HelpName}} myminio/mybucket
`,
}

// checkversionInfoSyntax - validate all the passed arguments
func checkEncryptInfoSyntax(ctx context.Context, cmd *cli.Command) {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

type encryptInfoMessage struct {
	Op         string `json:"op"`
	Status     string `json:"status"`
	URL        string `json:"url"`
	Encryption struct {
		Algorithm string `json:"algorithm,omitempty"`
		KeyID     string `json:"keyId,omitempty"`
	} `json:"encryption,omitempty"`
}

func (v encryptInfoMessage) JSON() string {
	v.Status = "success"
	jsonMessageBytes, e := json.MarshalIndent(v, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(jsonMessageBytes)
}

func (v encryptInfoMessage) String() string {
	msg := ""
	switch v.Encryption.Algorithm {
	case "":
		msg = fmt.Sprintf("Auto encryption is not enabled for %s ", v.URL)
	default:
		msg = "Auto encryption 'sse-s3' is enabled"
	}
	if v.Encryption.KeyID != "" {
		msg = fmt.Sprintf("Auto encryption 'sse-kms' is enabled with KeyID: %s", v.Encryption.KeyID)
	}
	return console.Colorize("encryptInfoMessage", msg)
}

func mainEncryptInfo(ctx context.Context, cmd *cli.Command) error {
	ctx, cancelEncryptInfo := context.WithCancel(globalContext)
	defer cancelEncryptInfo()

	console.SetColor("encryptInfoMessage", color.New(color.FgGreen))

	checkEncryptInfoSyntax(ctx, cmd)

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)
	// Create a new Client
	client, err := newClient(aliasedURL)
	fatalIf(err, "Unable to initialize connection.")
	algorithm, keyID, e := client.GetEncryption(ctx)
	fatalIf(e, "Unable to get encryption info")
	msg := encryptInfoMessage{
		Op:     "info",
		Status: "success",
		URL:    aliasedURL,
	}
	msg.Encryption.Algorithm = algorithm
	msg.Encryption.KeyID = keyID
	printMsg(msg)
	return nil
}
