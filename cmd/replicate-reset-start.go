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
	"strings"
	"time"

	"github.com/fatih/color"
	json "github.com/openstor/colorjson"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/openstor-go/v7/pkg/replication"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var replicateResyncStartFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "older-than",
		Usage: "replicate back objects older than value in duration string (e.g. 7d10h31s)",
	},
	&cli.StringFlag{
		Name:  "remote-bucket",
		Usage: "remote bucket ARN",
	},
}

var replicateResyncStartCmd = cli.Command{
	Name:         "start",
	Usage:        "start replicating back all previously replicated objects",
	Action:       mainReplicateResyncStart,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        append(globalFlags, replicateResyncStartFlags...),
	CustomHelpTemplate: `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{.HelpName}} TARGET

FLAGS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
EXAMPLES:
  1. Re-replicate previously replicated objects in bucket "mybucket" for alias "myminio" for remote target.
   {{.Prompt}} {{.HelpName}} myminio/mybucket --remote-bucket "arn:minio:replication::xxx:mybucket"

  2. Re-replicate all objects older than 60 days in bucket "mybucket" for remote bucket target.
   {{.Prompt}} {{.HelpName}} myminio/mybucket --older-than 60d --remote-bucket "arn:minio:replication::xxx:mybucket"
`,
}

// checkReplicateResyncStartSyntax - validate all the passed arguments
func checkReplicateResyncStartSyntax(ctx context.Context, cmd *cli.Command) {
	if cmd.Args().Len() != 1 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
	if cmd.String("remote-bucket") == "" {
		fatal(errDummy().Trace(), "--remote-bucket flag needs to be specified.")
	}
}

type replicateResyncMessage struct {
	Op                string                        `json:"op"`
	URL               string                        `json:"url"`
	ResyncTargetsInfo replication.ResyncTargetsInfo `json:"resyncInfo"`
	Status            string                        `json:"status"`
	TargetArn         string                        `json:"targetArn"`
}

func (r replicateResyncMessage) JSON() string {
	r.Status = "success"
	jsonMessageBytes, e := json.MarshalIndent(r, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(jsonMessageBytes)
}

func (r replicateResyncMessage) String() string {
	if len(r.ResyncTargetsInfo.Targets) == 1 {
		return console.Colorize("replicateResyncMessage", fmt.Sprintf("Replication reset started for %s with ID %s", r.URL, r.ResyncTargetsInfo.Targets[0].ResetID))
	}
	return console.Colorize("replicateResyncMessage", fmt.Sprintf("Replication reset started for %s", r.URL))
}

func mainReplicateResyncStart(ctx context.Context, cmd *cli.Command) error {
	ctx, cancelReplicateResyncStart := context.WithCancel(globalContext)
	defer cancelReplicateResyncStart()

	console.SetColor("replicateResyncMessage", color.New(color.FgGreen))

	checkReplicateResyncStartSyntax(ctx, cmd)

	// Get the alias parameter from cli
	args := cmd.Args()
	aliasedURL := args.Get(0)
	// Create a new Client
	client, err := newClient(aliasedURL)
	fatalIf(err, "Unable to initialize connection.")
	var olderThanStr string
	var olderThan time.Duration
	if cmd.IsSet("older-than") {
		olderThanStr = cmd.String("older-than")
		if olderThanStr != "" {
			days, e := ParseDuration(olderThanStr)
			if e != nil || !strings.ContainsAny(olderThanStr, "dwy") {
				fatalIf(probe.NewError(e), "Unable to parse older-than=`"+olderThanStr+"`.")
			}
			if days == 0 {
				fatalIf(probe.NewError(e), "older-than cannot be set to zero")
			}
			olderThan = time.Duration(days.Days())
		}
	}

	rinfo, err := client.ResetReplication(ctx, olderThan, cmd.String("remote-bucket"))
	fatalIf(err.Trace(args.Slice()...), "Unable to reset replication")
	printMsg(replicateResyncMessage{
		Op:                "start",
		URL:               aliasedURL,
		ResyncTargetsInfo: rinfo,
	})
	return nil
}
