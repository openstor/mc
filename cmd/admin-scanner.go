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

var adminScannerSubcommands = []cli.Command{
	adminScannerInfo,
	adminScannerTraceCmd,
}

var adminScannerCmd = cli.Command{
	Name:     "scanner",
	Usage:    "provide MinIO scanner info",
	Action:   mainAdminScanner,
	Before:   setGlobalsFromContext,
	Flags:    globalFlags,
	HideHelp: true,
}

// mainAdminScanner is the handle for "mc admin scanner" command.
func mainAdminScanner(ctx context.Context, cmd *cli.Command) error {
	commandNotFound(ctx, cmd, adminScannerSubcommands)
	return nil
}
