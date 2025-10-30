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

var quotaSubcommands = []*cli.Command{
	&quotaSetCmd,
	&quotaInfoCmd,
	&quotaClearCmd,
}

var quotaCmd = cli.Command{
	Name:            "quota",
	Usage:           "manage bucket quota",
	Action:          mainQuota,
	Before:          setGlobalsFromContext,
	Flags:           globalFlags,
	Commands:        quotaSubcommands,
	HideHelpCommand: true,
}

// mainQuota is the handle for "mc quota" command.
func mainQuota(ctx context.Context, cmd *cli.Command) error {
	// Convert []*cli.Command to []cli.Command for compatibility
	var subCmds []cli.Command
	for _, c := range quotaSubcommands {
		subCmds = append(subCmds, *c)
	}
	commandNotFound(ctx, cmd, subCmds)
	return nil
	// Sub-commands like "set", "clear", "info" have their own main.
}
