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
	"encoding/json"
	"os"
	"strings"

	"github.com/openstor/mc/pkg/probe"
	"github.com/urfave/cli/v3"
)

var aliasImportCmd = cli.Command{
	Name:            "import",
	Aliases:         []string{"i"},
	Usage:           "import configuration info to configuration file from a JSON formatted string ",
	Action:          mainAliasImport,
	OnUsageError:    onUsageError,
	Before:          setGlobalsFromContext,
	Flags:           globalFlags,
	HideHelpCommand: true,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} ALIAS ./credentials.json

  Credentials to be imported must be in the following JSON format:
  
  {
    "url": "http://localhost:9000",
    "accessKey": "YJ0RI0F4R5HWY38MD873",
    "secretKey": "OHz5CT7xdMHiXnKZP0BmZ5P4G5UvWvVaxR8gljLG",
    "api": "s3v4",
    "path": "auto"
  }

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Import the provided credentials.json file as 'myminio' to the config:
     {{ .Prompt }} {{ .HelpName }} myminio/ ./credentials.json

  2. Import the credentials through standard input as 'myminio' to the config:
     {{ .Prompt }} cat credentials.json | {{ .HelpName }} myminio/
`,
}

// checkAliasImportSyntax - verifies input arguments to 'alias import'.
func checkAliasImportSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()

	if cmd.NArg() == 0 {
		showCommandHelpAndExit(ctx, cmd, 1)
	}
	if cmd.NArg() > 2 {
		fatalIf(errInvalidArgument().Trace(cmd.Args().Tail()...),
			"Incorrect number of arguments for alias Import command.")
	}

	alias := cleanAlias(args.Get(0))

	if !isValidAlias(alias) {
		fatalIf(errInvalidAlias(alias), "Invalid alias.")
	}
}

func checkCredentialsSyntax(credentials aliasConfigV10) {
	if !isValidHostURL(credentials.URL) {
		fatalIf(errInvalidURL(credentials.URL), "Invalid URL.")
	}

	if !isValidAccessKey(credentials.AccessKey) {
		fatalIf(errInvalidArgument().Trace(credentials.AccessKey),
			"Invalid access key `"+credentials.AccessKey+"`.")
	}

	if !isValidSecretKey(credentials.SecretKey) {
		fatalIf(errInvalidArgument().Trace(),
			"Invalid secret key.")
	}

	if credentials.API != "" && !isValidAPI(credentials.API) { // Empty value set to default "S3v4".
		fatalIf(errInvalidArgument().Trace(credentials.API),
			"Unrecognized API signature. Valid options are `[S3v4, S3v2]`.")
	}
	if !isValidPath(credentials.Path) {
		fatalIf(errInvalidArgument().Trace(credentials.Path),
			"Unrecognized path value. Valid options are `[auto, on, off]`.")
	}
}

// importAlias - set an alias config based on imported values.
func importAlias(alias string, aliasCfgV10 aliasConfigV10) aliasMessage {
	checkCredentialsSyntax(aliasCfgV10)

	mcCfgV10, err := loadMcConfig()
	fatalIf(err.Trace(globalMCConfigVersion), "Unable to load config `"+mustGetMcConfigPath()+"`.")

	// Add new host.
	mcCfgV10.Aliases[alias] = aliasCfgV10
	fatalIf(saveMcConfig(mcCfgV10).Trace(alias), "Unable to import credentials to `"+mustGetMcConfigPath()+"`.")
	return aliasMessage{
		Alias:     alias,
		URL:       mcCfgV10.Aliases[alias].URL,
		AccessKey: mcCfgV10.Aliases[alias].AccessKey,
		SecretKey: mcCfgV10.Aliases[alias].SecretKey,
		API:       mcCfgV10.Aliases[alias].API,
		Path:      mcCfgV10.Aliases[alias].Path,
	}
}

func mainAliasImport(ctx context.Context, cmd *cli.Command) error {
	var (
		args  = cmd.Args()
		alias = cleanAlias(args.Get(0))
	)

	checkAliasImportSyntax(ctx, cmd)
	var credentialsJSON aliasConfigV10

	credsFile := strings.TrimSpace(args.Get(1))
	if credsFile == "" {
		credsFile = os.Stdin.Name()
	}
	input, e := os.ReadFile(credsFile)
	fatalIf(probe.NewError(e).Trace(credsFile), "Unable to parse credentials file")

	e = json.Unmarshal(input, &credentialsJSON)
	fatalIf(probe.NewError(e).Trace(credsFile), "Unable to parse input credentials")

	msg := importAlias(alias, credentialsJSON)
	msg.op = cmd.Name

	printMsg(msg)

	return nil
}
