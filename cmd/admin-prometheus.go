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

var adminPrometheusSubcommands = []cli.Command{
	adminPrometheusGenerateCmd,
	adminPrometheusMetricsCmd,
}

var adminPrometheusCmd = cli.Command{
	Name:     "prometheus",
	Usage:    "manages prometheus config",
	Action:   mainAdminPrometheus,
	Before:   setGlobalsFromContext,
	Flags:    globalFlags,
	HideHelp: true,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
`,
}

// mainAdminPrometheus is the handle for "mc admin prometheus" command.
func mainAdminPrometheus(ctx context.Context, cmd *cli.Command) error {
	commandNotFound(ctx, cmd, adminPrometheusSubcommands)
	return nil
}
