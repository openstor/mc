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

var adminKMSKeyStatusCmd = &cli.Command{
	Name:         "status",
	Usage:        "request status information for a KMS master key",
	Action:       mainAdminKMSKeyStatus,
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
  1. Get default master key and its status from a MinIO server/cluster.
     $ {{.HelpName}} play
  2. Get the status of one particular master key from a MinIO server/cluster.
     $ {{.HelpName}} play my-master-key
`,
}

// adminKMSKeyCmd is the handle for the "mc admin kms key" command.
func mainAdminKMSKeyStatus(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	if args.Len() == 0 || args.Len() > 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}

	console.SetColor("StatusSuccess", color.New(color.FgGreen, color.Bold))
	console.SetColor("StatusError", color.New(color.FgRed, color.Bold))
	console.SetColor("StatusUnknown", color.New(color.FgYellow, color.Bold))

	client, err := newAdminClient(args.Get(0))
	fatalIf(err, "Unable to get a configured admin connection.")

	var keyID string
	if args.Len() == 2 {
		keyID = args.Get(1)
	}
	status, e := client.GetKeyStatus(globalContext, keyID)
	fatalIf(probe.NewError(e), "Failed to get status information")

	printMsg(kmsKeyStatusMsg{
		KeyID:         status.KeyID,
		EncryptionErr: status.EncryptionErr,
		DecryptionErr: status.DecryptionErr,
	})
	return nil
}

type kmsKeyStatusMsg struct {
	KeyID         string `json:"keyId"`
	EncryptionErr string `json:"encryptionError,omitempty"`
	DecryptionErr string `json:"decryptionError,omitempty"`
	Status        string `json:"status"`
}

func (s kmsKeyStatusMsg) JSON() string {
	s.Status = "success"
	kmsBytes, e := json.MarshalIndent(s, "", "    ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(kmsBytes)
}

func (s kmsKeyStatusMsg) String() string {
	msg := fmt.Sprintf("Key: %s\n", s.KeyID)

	success := console.Colorize("StatusSuccess", "✔")
	failure := console.Colorize("StatusError", "✗")
	dunno := console.Colorize("StatusUnknown", "?")

	formatStatus := func(name string, unknown bool, err string) string {
		st := ""
		switch {
		case !unknown && err == "":
			st = success
		case unknown:
			st = dunno
		case err != "":
			st = fmt.Sprintf("%s (%s)", failure, err)
		}
		return fmt.Sprintf("   - %s %s\n", name, st)
	}

	msg += formatStatus("Encryption", false, s.EncryptionErr)
	msg += formatStatus("Decryption", s.EncryptionErr != "", s.DecryptionErr)
	return msg
}
