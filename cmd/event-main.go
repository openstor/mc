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

	"github.com/urfave/cli/v3"
)

var eventFlags = []cli.Flag{}

var eventSubcommands = []*cli.Command{
	&eventAddCmd,
	&eventRemoveCmd,
	&eventListCmd,
}

var eventCmd = cli.Command{
	Name:            "event",
	Usage:           "manage object notifications",
	HideHelpCommand: true,
	Action:          mainEvent,
	Before:          setGlobalsFromContext,
	Flags:           append(eventFlags, globalFlags...),
	Commands:        eventSubcommands,
}

// mainEvent is the handle for "mc event" command.
func mainEvent(ctx context.Context, cmd *cli.Command) error {
	var cmds []cli.Command
	for _, c := range eventSubcommands {
		cmds = append(cmds, *c)
	}
	commandNotFound(ctx, cmd, cmds)
	return nil
	// Sub-commands like "add", "remove", "list" have their own main.
}
