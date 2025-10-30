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

var adminServiceSubcommands = []*cli.Command{
	adminServiceRestartCmd,
	adminServiceStopCmd,
	adminServiceUnfreezeCmd,
	adminServiceFreezeCmd,
}

var adminServiceCmd = cli.Command{
	Name:  "service",
	Usage: "restart or unfreeze a MinIO cluster",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		// Convert []*cli.Command to []cli.Command
		var commands []cli.Command
		for _, c := range adminServiceSubcommands {
			commands = append(commands, *c)
		}
		commandNotFound(ctx, cmd, commands)
		return nil
	},
	Before:          setGlobalsFromContext,
	Flags:           globalFlags,
	HideHelpCommand: true,
	Commands:        adminServiceSubcommands,
}

// mainAdmin is the handle for "mc admin service" command.
func mainAdminService(ctx context.Context, cmd *cli.Command) error {
	// Convert []*cli.Command to []cli.Command
	var commands []cli.Command
	for _, c := range adminServiceSubcommands {
		commands = append(commands, *c)
	}
	commandNotFound(ctx, cmd, commands)
	return nil
	// Sub-commands like "status", "restart" have their own main.
}
