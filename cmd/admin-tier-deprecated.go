// Copyright (c) 2022 MinIO, Inc.
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

// Wrapper functions for deprecated tier commands
func mainAdminTierInfoWrapper(ctx context.Context, cmd *cli.Command) error {
	return mainAdminTierInfo(ctx, cmd)
}

func mainAdminTierListDeprecated(ctx context.Context, cmd *cli.Command) error {
	return mainAdminTierList(ctx, cmd)
}

func mainAdminTierAddDeprecated(ctx context.Context, cmd *cli.Command) error {
	return mainAdminTierAdd(ctx, cmd)
}

func mainAdminTierEditDeprecated(ctx context.Context, cmd *cli.Command) error {
	return mainAdminTierEdit(ctx, cmd)
}

func mainAdminTierVerifyDeprecated(ctx context.Context, cmd *cli.Command) error {
	return mainAdminTierVerify(ctx, cmd)
}

func mainAdminTierRmDeprecated(ctx context.Context, cmd *cli.Command) error {
	return mainAdminTierRm(ctx, cmd)
}

var adminTierDepCmds = []cli.Command{
	adminTierDepInfoCmd,
	adminTierDepListCmd,
	adminTierDepAddCmd,
	adminTierDepEditCmd,
	adminTierDepVerifyCmd,
	adminTierDepRmCmd,
}

var (
	adminTierDepInfoCmd = cli.Command{
		Name:         "info",
		Usage:        "display tier statistics",
		Action:       mainAdminTierInfoWrapper,
		Hidden:       true,
		OnUsageError: onUsageError,
		Before:       setGlobalsFromContext,
		Flags:        globalFlags,
		CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} ALIAS [NAME]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}

EXAMPLES:
  1. Prints per-tier statistics of all remote tier targets configured on 'myminio':
     {{.Prompt}} {{.HelpName}} myminio

  2. Print per-tier statistics of given tier name 'MINIOTIER-1':
     {{.Prompt}} {{.HelpName}} myminio MINIOTIER-1
`,
	}

	adminTierDepListCmd = cli.Command{
		Name:         "ls",
		Usage:        "lists configured remote tier targets",
		Action:       mainAdminTierListDeprecated,
		Hidden:       true,
		OnUsageError: onUsageError,
		Before:       setGlobalsFromContext,
		Flags:        globalFlags,
		CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} ALIAS

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}

EXAMPLES:
  1. List remote tier targets configured on 'myminio':
     {{.Prompt}} {{.HelpName}} myminio
`,
	}

	adminTierDepAddCmd = cli.Command{
		Name:         "add",
		Usage:        "add a new remote tier target",
		Action:       mainAdminTierAddDeprecated,
		Hidden:       true,
		OnUsageError: onUsageError,
		Before:       setGlobalsFromContext,
		Flags:        append(globalFlags, adminTierAddFlags...),
		CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TYPE ALIAS NAME [FLAGS]

TYPE:
  Type of the cloud storage backend to add. Supported values are minio, s3, azure and gcs.

NAME:
  Name of the remote tier target. e.g WARM-TIER

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Configure a new remote tier which transitions objects to a bucket in a MinIO deployment:
     {{.Prompt}} {{.HelpName}} minio myminio WARM-MINIO-TIER --endpoint https://warm-minio.com \
        --access-key ACCESSKEY --secret-key SECRETKEY --bucket mybucket --prefix myprefix/

  2. Configure a new remote tier which transitions objects to a bucket in Azure Blob Storage:
     {{.Prompt}} {{.HelpName}} azure myminio AZTIER --account-name ACCOUNT-NAME --account-key ACCOUNT-KEY \
        --bucket myazurebucket --prefix myazureprefix/

  3. Configure a new remote tier which transitions objects to a bucket in AWS S3 with STANDARD storage class:
     {{.Prompt}} {{.HelpName}} s3 myminio S3TIER --endpoint https://s3.amazonaws.com \
        --access-key ACCESSKEY --secret-key SECRETKEY --bucket mys3bucket --prefix mys3prefix/ \
        --storage-class "STANDARD" --region us-west-2

  4. Configure a new remote tier which transitions objects to a bucket in Google Cloud Storage:
     {{.Prompt}} {{.HelpName}} gcs myminio GCSTIER --credentials-file /path/to/credentials.json \
        --bucket mygcsbucket  --prefix mygcsprefix/
`,
	}
	adminTierDepEditCmd = cli.Command{
		Name:         "edit",
		Usage:        "update an existing remote tier configuration",
		Action:       mainAdminTierEditDeprecated,
		Hidden:       true,
		OnUsageError: onUsageError,
		Before:       setGlobalsFromContext,
		Flags:        append(globalFlags, adminTierEditFlags...),
		CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} ALIAS NAME

NAME:
  Name of remote tier. e.g WARM-TIER

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Update credentials for an existing Azure Blob Storage remote tier:
     {{.Prompt}} {{.HelpName}} myminio AZTIER --account-key ACCOUNT-KEY

  2. Update credentials for an existing AWS S3 compatible remote tier:
     {{.Prompt}} {{.HelpName}} myminio S3TIER --access-key ACCESS-KEY --secret-key SECRET-KEY

  3. Update credentials for an existing Google Cloud Storage remote tier:
     {{.Prompt}} {{.HelpName}} myminio GCSTIER --credentials-file /path/to/credentials.json
`,
	}

	adminTierDepVerifyCmd = cli.Command{
		Name:         "verify",
		Usage:        "verifies if remote tier configuration is valid",
		Action:       mainAdminTierVerifyDeprecated,
		Hidden:       true,
		OnUsageError: onUsageError,
		Before:       setGlobalsFromContext,
		Flags:        globalFlags,
		CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET NAME

NAME:
  Name of remote tier target. e.g WARM-TIER

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Verify if a tier config is valid.
     {{.Prompt}} {{.HelpName}} myminio WARM-TIER
`,
	}

	adminTierDepRmCmd = cli.Command{
		Name:         "rm",
		Usage:        "removes an empty remote tier",
		Action:       mainAdminTierRmDeprecated,
		Hidden:       true,
		OnUsageError: onUsageError,
		Before:       setGlobalsFromContext,
		Flags:        globalFlags,
		CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} ALIAS NAME

NAME:
  Name of remote tier target. e.g WARM-TIER

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Remove an empty tier by name 'WARM-TIER':
     {{.Prompt}} {{.HelpName}} myminio WARM-TIER
`,
	}
)
