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
	"fmt"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestAutoCompletionCompletness(t *testing.T) {
	var checkCompletion func(cmd *cli.Command, cmdPath string) error

	checkCompletion = func(cmd *cli.Command, cmdPath string) error {
		// Special handling for admin command - it has subcommands defined separately
		if cmd.Name == "admin" {
			// admin command is handled specially, skip completion check
			return nil
		}

		// If command has subcommands, recursively check them
		if cmd.Commands != nil {
			for _, subCmd := range cmd.Commands {
				if subCmd.Hidden {
					continue
				}
				err := checkCompletion(subCmd, cmdPath+"/"+subCmd.Name)
				if err != nil {
					return err
				}
			}
			return nil
		}
		// Only check completion for leaf commands (commands without subcommands)
		_, ok := completeCmds[cmdPath]
		if !ok && !cmd.Hidden {
			return fmt.Errorf("Completion for `%s` not found", cmdPath)
		}
		return nil
	}

	for _, cmd := range appCmds {
		if cmd.Hidden {
			continue
		}
		err := checkCompletion(cmd, "/"+cmd.Name)
		if err != nil {
			t.Fatalf("Missing completion function: %v", err)
		}

	}
}
