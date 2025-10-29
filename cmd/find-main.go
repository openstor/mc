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
	"regexp"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
)

// List of all flags supported by find command.
var (
	findFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "exec",
			Usage: "spawn an external process for each matching object (see FORMAT)",
		},
		&cli.StringFlag{
			Name:  "ignore",
			Usage: "exclude objects matching the wildcard pattern",
		},
		&cli.BoolFlag{
			Name:  "versions",
			Usage: "include all objects versions",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "find object names matching wildcard pattern",
		},
		&cli.StringFlag{
			Name:  "newer-than",
			Usage: "match all objects newer than value in duration string (e.g. 7d10h31s)",
		},
		&cli.StringFlag{
			Name:  "older-than",
			Usage: "match all objects older than value in duration string (e.g. 7d10h31s)",
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "match directory names matching wildcard pattern",
		},
		&cli.StringFlag{
			Name:  "print",
			Usage: "print in custom format to STDOUT (see FORMAT)",
		},
		&cli.StringFlag{
			Name:  "regex",
			Usage: "match directory and object name with RE2 regex pattern",
		},
		&cli.StringFlag{
			Name:  "larger",
			Usage: "match all objects larger than specified size in units (see UNITS)",
		},
		&cli.StringFlag{
			Name:  "smaller",
			Usage: "match all objects smaller than specified size in units (see UNITS)",
		},
		&cli.UintFlag{
			Name:  "maxdepth",
			Usage: "limit directory navigation to specified depth",
		},
		&cli.BoolFlag{
			Name:  "watch",
			Usage: "monitor a specified path for newly created object(s)",
		},
		&cli.StringSliceFlag{
			Name:  "metadata",
			Usage: "match metadata with RE2 regex pattern. Specify each with key=regex. MinIO server only.",
		},
		&cli.StringSliceFlag{
			Name:  "tags",
			Usage: "match tags with RE2 regex pattern. Specify each with key=regex. MinIO server only.",
		},
	}
)

var findCmd = cli.Command{
	Name:  "find",
	Usage: "search for objects",
	Action: func(ctx context.Context, cliCtx *cli.Command) error {
		return mainFind(cliCtx)
	},
	Before: setGlobalsFromContext,
	Flags:  append(findFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Find all "go" files in mybucket.
     {{.Prompt}} {{.HelpName}} s3/mybucket --name "*.go"

  2. Find all objects in mybucket that were created within the last 7 days.
     {{.Prompt}} {{.HelpName}} s3/mybucket --created-within 7d

  3. Find all objects in mybucket that were created more than 7 days ago.
     {{.Prompt}} {{.HelpName}} s3/mybucket --created-before 7d

  4. Find all objects in mybucket that are larger than 1MB.
     {{.Prompt}} {{.HelpName}} s3/mybucket --larger 1MB

  5. Find all objects in mybucket that are smaller than 1MB.
     {{.Prompt}} {{.HelpName}} s3/mybucket --smaller 1MB

  6. Find all objects in mybucket that have a custom metadata field "x-amz-meta-author" with the value "john".
     {{.Prompt}} {{.HelpName}} s3/mybucket --metadata "x-amz-meta-author=john"

  7. Find all objects in mybucket that have a custom metadata field "x-amz-meta-author" with any value.
     {{.Prompt}} {{.HelpName}} s3/mybucket --metadata "x-amz-meta-author"

  8. Find all objects in mybucket that have a custom tag "author" with the value "john".
     {{.Prompt}} {{.HelpName}} s3/mybucket --tags "author=john"

  9. Find all objects in mybucket that have a custom tag "author" with any value.
     {{.Prompt}} {{.HelpName}} s3/mybucket --tags "author"

  10. Find all objects in mybucket that match the specified query.
      {{.Prompt}} {{.HelpName}} s3/mybucket --query "select * from s3object where author = 'john'"

  11. Find all objects in mybucket and count them.
      {{.Prompt}} {{.HelpName}} s3/mybucket --count

  12. Find all objects in mybucket and print them in a tree format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --tree

  13. Find all objects in mybucket and print them in a table format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --table

  14. Find all objects in mybucket and print them in a json format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --json

  15. Find all objects in mybucket and print them in a csv format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --csv

  16. Find all objects in mybucket and print them in a yaml format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --yaml

  17. Find all objects in mybucket and print them in a xml format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --xml

  18. Find all objects in mybucket and print them in a html format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --html

  19. Find all objects in mybucket and print them in a markdown format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --markdown

  20. Find all objects in mybucket and print them in a custom format.
      {{.Prompt}} {{.HelpName}} s3/mybucket --format "{{.Name}} {{.Size}} {{.ModTime}}"
`,
}

// checkFindSyntax - validate the passed arguments
func checkFindSyntax(ctx context.Context, cliCtx *cli.Command, encKeyDB map[string][]prefixSSEPair) {
	args := cliCtx.Args().Slice()
	if len(args) == 0 {
		args = []string{"./"} // No args just default to present directory.
	} else if args[0] == "." {
		args[0] = "./" // If the arg is '.' treat it as './'.
	}

	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			fatalIf(errInvalidArgument().Trace(args...), "Unable to validate empty argument.")
		}
	}

	// Extract input URLs and validate.
	for _, url := range args {
		_, _, err := url2Stat(ctx, url2StatOptions{urlStr: url, versionID: "", fileAttr: false, encKeyDB: encKeyDB, timeRef: time.Time{}, isZip: false, ignoreBucketExistsCheck: false})
		if err != nil {
			// Bucket name empty is a valid error for 'find myminio' unless we are using watch, treat it as such.
			if _, ok := err.ToGoError().(BucketNameEmpty); ok && !cliCtx.Bool("watch") {
				continue
			}
			fatalIf(err.Trace(url), "Unable to stat `"+url+"`.")
		}
	}
}

// Find context is container to hold all parsed input arguments,
// each parsed input is stored in its native typed form for
// ease of repurposing.
type findContext struct {
	*cli.Command
	execCmd       string
	ignorePattern string
	namePattern   string
	pathPattern   string
	regexPattern  *regexp.Regexp
	maxDepth      uint
	printFmt      string
	olderThan     string
	newerThan     string
	largerSize    uint64
	smallerSize   uint64
	watch         bool
	withVersions  bool
	matchMeta     map[string]*regexp.Regexp
	matchTags     map[string]*regexp.Regexp

	// Internal values
	targetAlias   string
	targetURL     string
	targetFullURL string
	clnt          Client
}

// mainFind - handler for mc find commands
func mainFind(cliCtx *cli.Command) error {
	ctx, cancelFind := context.WithCancel(globalContext)
	defer cancelFind()

	// Additional command specific theme customization.
	console.SetColor("Find", color.New(color.FgGreen, color.Bold))
	console.SetColor("FindExecErr", color.New(color.FgRed, color.Italic, color.Bold))

	// Parse encryption keys per command.
	encKeyDB, err := validateAndCreateEncryptionKeys(ctx, cliCtx)
	fatalIf(err, "Unable to parse encryption keys.")

	checkFindSyntax(ctx, cliCtx, encKeyDB)

	args := cliCtx.Args().Slice()
	if len(args) == 0 {
		args = []string{"./"} // Not args present default to present directory.
	} else if args[0] == "." {
		args[0] = "./" // If the arg is '.' treat it as './'.
	}

	clnt, err := newClient(args[0])
	fatalIf(err.Trace(args...), "Unable to initialize `"+args[0]+"`.")

	var olderThan, newerThan string

	if cliCtx.String("older-than") != "" {
		olderThan = cliCtx.String("older-than")
	}
	if cliCtx.String("newer-than") != "" {
		newerThan = cliCtx.String("newer-than")
	}

	// Use 'e' to indicate Go error, this is a convention followed in `mc`. For probe.Error we call it
	// 'err' and regular Go error is called as 'e'.
	var e error
	var largerSize, smallerSize uint64

	if cliCtx.String("larger") != "" {
		largerSize, e = humanize.ParseBytes(cliCtx.String("larger"))
		fatalIf(probe.NewError(e).Trace(cliCtx.String("larger")), "Unable to parse input bytes.")
	}

	if cliCtx.String("smaller") != "" {
		smallerSize, e = humanize.ParseBytes(cliCtx.String("smaller"))
		fatalIf(probe.NewError(e).Trace(cliCtx.String("smaller")), "Unable to parse input bytes.")
	}

	// Get --versions flag
	withVersions := cliCtx.Bool("versions")

	targetAlias, _, hostCfg, err := expandAlias(args[0])
	fatalIf(err.Trace(args[0]), "Unable to expand alias.")

	var targetFullURL string
	if hostCfg != nil {
		targetFullURL = hostCfg.URL
	}
	var regMatch *regexp.Regexp
	if cliCtx.String("regex") != "" {
		regMatch = regexp.MustCompile(cliCtx.String("regex"))
	}

	return doFind(ctx, &findContext{
		Command:       cliCtx,
		maxDepth:      cliCtx.Uint("maxdepth"),
		execCmd:       cliCtx.String("exec"),
		printFmt:      cliCtx.String("print"),
		namePattern:   cliCtx.String("name"),
		pathPattern:   cliCtx.String("path"),
		regexPattern:  regMatch,
		ignorePattern: cliCtx.String("ignore"),
		withVersions:  withVersions,
		olderThan:     olderThan,
		newerThan:     newerThan,
		largerSize:    largerSize,
		smallerSize:   smallerSize,
		watch:         cliCtx.Bool("watch"),
		targetAlias:   targetAlias,
		targetURL:     args[0],
		targetFullURL: targetFullURL,
		clnt:          clnt,
		matchMeta:     getRegexMap(ctx, cliCtx, "metadata"),
		matchTags:     getRegexMap(ctx, cliCtx, "tags"),
	})
}
