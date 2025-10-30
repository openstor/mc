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

var adminBucketRemoteSubcommands = []*cli.Command{
	&adminBucketRemoteAddCmd,
	&adminBucketRemoteEditCmd,
	&adminBucketRemoteRmCmd,
}

var adminBucketRemoteCmd = cli.Command{
	Name:            "remote",
	Usage:           "configure remote target buckets",
	Action:          mainadminBucketRemote,
	Before:          setGlobalsFromContext,
	Flags:           globalFlags,
	Commands:        adminBucketRemoteSubcommands,
	HideHelpCommand: true,
}

// mainadminBucketRemote is the handle for "mc admin bucket remote" command.
func mainadminBucketRemote(ctx context.Context, cmd *cli.Command) error {
	var cmds []cli.Command
	for _, c := range adminBucketRemoteSubcommands {
		cmds = append(cmds, *c)
	}
	commandNotFound(ctx, cmd, cmds)
	return nil
	// Sub-commands like "add", "ls", "rm" have their own main.
}
