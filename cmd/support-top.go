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

	"github.com/urfave/cli/v3"
)

var supportTopSubcommands = []*cli.Command{
	&supportTopAPICmd,
	&supportTopDriveCmd,
	&supportTopLocksCmd,
	&supportTopNetCmd,
	&supportTopRPCCmd,
}

var supportTopCmd = cli.Command{
	Name:            "top",
	Usage:           "provide top like statistics for MinIO",
	Action:          mainSupportTop,
	Before:          setGlobalsFromContext,
	Flags:           globalFlags,
	Commands:        supportTopSubcommands,
	HideHelpCommand: true,
}

// mainSupportTop is the handle for "mc support top" command.
func mainSupportTop(ctx context.Context, cmd *cli.Command) error {
	// Convert []*cli.Command to []cli.Command for compatibility
	var subCmds []cli.Command
	for _, c := range supportTopSubcommands {
		subCmds = append(subCmds, *c)
	}
	commandNotFound(ctx, cmd, subCmds)
	return nil
	// Sub-commands like "locks" have their own main.
}
