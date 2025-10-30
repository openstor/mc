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
	"errors"

	"github.com/fatih/color"
	json "github.com/openstor/colorjson"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

var eventRemoveFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "force",
		Usage: "force removing all bucket notifications",
	},
	&cli.StringFlag{
		Name:  "event",
		Value: "put,delete,get",
		Usage: "filter specific type of event. Defaults to all event",
	},
	&cli.StringFlag{
		Name:  "prefix",
		Usage: "filter event associated to the specified prefix",
	},
	&cli.StringFlag{
		Name:  "suffix",
		Usage: "filter event associated to the specified suffix",
	},
}

var eventRemoveCmd = cli.Command{
	Name:         "remove",
	Aliases:      []string{"rm"},
	Usage:        "remove a bucket notification; '--force' removes all bucket notifications",
	Action:       mainEventRemove,
	OnUsageError: onUsageError,
	Before:       setGlobalsFromContext,
	Flags:        append(eventRemoveFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [ARN] [FLAGS]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Remove bucket notification associated to a specific arn
    {{.Prompt}} {{.HelpName}} myminio/mybucket arn:aws:sqs:us-west-2:444455556666:your-queue

  2. Remove all bucket notifications. --force flag is mandatory here
    {{.Prompt}} {{.HelpName}} myminio/mybucket --force
`,
}

// checkEventRemoveSyntax - validate all the passed arguments
func checkEventRemoveSyntax(ctx context.Context, cmd *cli.Command) {
	if cmd.Args().Len() == 0 || cmd.Args().Len() > 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
	if cmd.Args().Len() == 1 && !cmd.Bool("force") {
		fatalIf(probe.NewError(errors.New("")), "--force flag needs to be passed to remove all bucket notifications.")
	}
}

// eventRemoveMessage container
type eventRemoveMessage struct {
	ARN    string `json:"arn"`
	Status string `json:"status"`
}

// JSON jsonified remove message.
func (u eventRemoveMessage) JSON() string {
	u.Status = "success"
	eventRemoveMessageJSONBytes, e := json.MarshalIndent(u, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(eventRemoveMessageJSONBytes)
}

func (u eventRemoveMessage) String() string {
	msg := console.Colorize("Event", "Successfully removed "+u.ARN)
	return msg
}

func mainEventRemove(ctx context.Context, cmd *cli.Command) error {
	ctx, cancelEventRemove := context.WithCancel(globalContext)
	defer cancelEventRemove()

	console.SetColor("Event", color.New(color.FgGreen, color.Bold))

	checkEventRemoveSyntax(ctx, cmd)

	args := cmd.Args()
	path := args.Get(0)

	arn := ""
	if args.Len() == 2 {
		arn = args.Get(1)
	}

	client, err := newClient(path)
	if err != nil {
		fatalIf(err.Trace(), "Unable to parse the provided url.")
	}

	s3Client, ok := client.(*S3Client)
	if !ok {
		fatalIf(errDummy().Trace(), "The provided url doesn't point to a S3 server.")
	}

	// flags for the attributes of the even
	event := cmd.String("event")
	prefix := cmd.String("prefix")
	suffix := cmd.String("suffix")

	err = s3Client.RemoveNotificationConfig(ctx, arn, event, prefix, suffix)
	if err != nil {
		fatalIf(err, "Unable to disable notification on the specified bucket.")
	}

	printMsg(eventRemoveMessage{ARN: arn})

	return nil
}
