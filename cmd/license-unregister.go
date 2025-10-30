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

const licUnregisterMsgTag = "licenseUnregisterMessage"

var licenseUnregisterCmd = cli.Command{
	Name:         "unregister",
	Usage:        "unregister from MinIO Subscription Network",
	OnUsageError: onUsageError,
	Action:       mainLicenseUnregister,
	Before:       setGlobalsFromContext,
	Hidden:       true,
	Flags:        subnetCommonFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Unregister MinIO cluster at alias 'myminio' from SUBNET
     {{.Prompt}} {{.HelpName}} myminio
`,
}

type licUnregisterMessage struct {
	Status string `json:"status"`
	Alias  string `json:"-"`
}

// String colorized license unregister message
func (li licUnregisterMessage) String() string {
	msg := fmt.Sprintf("%s unregistered successfully.", li.Alias)
	return console.Colorize(licUnregisterMsgTag, msg)
}

// JSON jsonified license unregister message
func (li licUnregisterMessage) JSON() string {
	jsonBytes, e := json.MarshalIndent(li, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(jsonBytes)
}

// checkLicenseUnregisterSyntax - validate arguments passed by a user
func checkLicenseUnregisterSyntax(ctx context.Context, cmd *cli.Command) {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

func mainLicenseUnregister(ctx context.Context, cmd *cli.Command) error {
	console.SetColor(licUnregisterMsgTag, color.New(color.FgGreen, color.Bold))
	checkLicenseUnregisterSyntax(ctx, cmd)

	aliasedURL := cmd.Args().Get(0)
	alias, apiKey := initSubnetConnectivity(ctx, cmd, aliasedURL, true)
	if len(apiKey) == 0 {
		// api key not passed as flag. Check that the cluster is registered.
		apiKey = validateClusterRegistered(alias, true)
	}

	if !globalAirgapped {
		info := getAdminInfo(aliasedURL)
		e := unregisterClusterFromSubnet(info.DeploymentID, apiKey)
		fatalIf(probe.NewError(e), "Could not unregister cluster from SUBNET:")
	}

	removeSubnetAuthConfig(alias)

	printMsg(licUnregisterMessage{Status: "success", Alias: alias})
	return nil
}
