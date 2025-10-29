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
	"runtime"

	"github.com/urfave/cli/v3"
)

func checkCopySyntax(ctx context.Context, cmd *cli.Command) {
	if cmd.NArg() < 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code.
	}
	parseChecksum(cmd)

	// extract URLs.
	args := cmd.Args()
	if args.Len() < 2 {
		fatalIf(errDummy().Trace(args.Tail()...), "Unable to parse source and target arguments.")
	}

	srcURLs := make([]string, args.Len()-1)
	for i := 0; i < args.Len()-1; i++ {
		srcURLs[i] = args.Get(i)
	}
	tgtURL := args.Get(args.Len() - 1)
	isZip := cmd.Bool("zip")
	versionID := cmd.String("version-id")

	if versionID != "" && len(srcURLs) > 1 {
		fatalIf(errDummy().Trace(args.Tail()...), "Unable to pass --version flag with multiple copy sources arguments.")
	}

	if isZip && cmd.String("rewind") != "" {
		fatalIf(errDummy().Trace(args.Tail()...), "--zip and --rewind cannot be used together")
	}

	// Check if bucket name is passed for URL type arguments.
	url := newClientURL(tgtURL)
	if url.Host != "" {
		if url.Path == string(url.Separator) {
			fatalIf(errInvalidArgument().Trace(), fmt.Sprintf("Target `%s` does not contain bucket name.", tgtURL))
		}
	}

	if cmd.String(rdFlag) != "" && cmd.String(rmFlag) == "" {
		fatalIf(errInvalidArgument().Trace(), fmt.Sprintf("Both object retention flags `--%s` and `--%s` are required.\n", rdFlag, rmFlag))
	}

	if cmd.String(rdFlag) == "" && cmd.String(rmFlag) != "" {
		fatalIf(errInvalidArgument().Trace(), fmt.Sprintf("Both object retention flags `--%s` and `--%s` are required.\n", rdFlag, rmFlag))
	}

	// Preserve functionality not supported for windows
	if cmd.Bool("preserve") && runtime.GOOS == "windows" {
		fatalIf(errInvalidArgument().Trace(), "Permissions are not preserved on windows platform.")
	}
}
