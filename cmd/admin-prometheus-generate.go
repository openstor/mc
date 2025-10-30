// Copyright (c) 2015-2024 MinIO, Inc.
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
	"net/url"
	"time"

	"github.com/fatih/color"
	"github.com/openstor/mc/pkg/probe"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"

	json "github.com/openstor/colorjson"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultJobName    = "minio-job"
	metricsV2BasePath = "/minio/v2/metrics"
)

var prometheusFlags = append(metricsFlags,
	&cli.BoolFlag{
		Name:  "public",
		Usage: "disable bearer token generation for scrape_configs",
	})

var adminPrometheusGenerateCmd = cli.Command{
	Name:            "generate",
	Usage:           "generates prometheus config",
	Action:          mainAdminPrometheusGenerate,
	OnUsageError:    onUsageError,
	Before:          setGlobalsFromContext,
	Flags:           append(prometheusFlags, globalFlags...),
	HideHelpCommand: true,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET [METRIC-TYPE]

METRIC-TYPE:
  valid values are
    api-version v2 ['cluster', 'node', 'bucket', 'resource']. defaults to 'cluster' if not specified.
    api-version v3 ["api", "system", "debug", "cluster", "ilm", "audit", "logger", "replication", "notification", "scanner"]. defaults to all if not specified.

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES (v3):
  1. Generate a default prometheus config.
     {{.Prompt}} {{.HelpName}} play --api-version v3

  2. Generate prometheus config for api metrics.
     {{.Prompt}} {{.HelpName}} play api --api-version v3

  3. Generate prometheus config for api metrics of bucket 'mybucket'.
     {{.Prompt}} {{.HelpName}} play api --bucket mybucket --api-version v3

  4. Generate prometheus config for system metrics.
     {{.Prompt}} {{.HelpName}} play system --api-version v3

  5. Generate prometheus config for debug metrics.
     {{.Prompt}} {{.HelpName}} play debug --api-version v3

  6. Generate prometheus config for cluster metrics.
     {{.Prompt}} {{.HelpName}} play cluster --api-version v3

  7. Generate prometheus config for ilm metrics.
     {{.Prompt}} {{.HelpName}} play ilm --api-version v3

  8. Generate prometheus config for audit metrics.
     {{.Prompt}} {{.HelpName}} play audit --api-version v3

  9. Generate prometheus config for logger metrics.
     {{.Prompt}} {{.HelpName}} play logger --api-version v3

  10. Generate prometheus config for replication metrics.
     {{.Prompt}} {{.HelpName}} play replication --api-version v3

  11. Generate prometheus config for replication metrics of bucket 'mybucket'.
     {{.Prompt}} {{.HelpName}} play replication --bucket mybucket --api-version v3

  12. Generate prometheus config for notification metrics.
     {{.Prompt}} {{.HelpName}} play notification --api-version v3

  13. Generate prometheus config for scanner metrics.
     {{.Prompt}} {{.HelpName}} play scanner --api-version v3

EXAMPLES (v2):
  1. Generate a default prometheus config.
     {{.Prompt}} {{.HelpName}} play

  2. Generate prometheus config for node metrics.
     {{.Prompt}} {{.HelpName}} play node

  3. Generate prometheus config for bucket metrics.
     {{.Prompt}} {{.HelpName}} play bucket

  4. Generate prometheus config for resource metrics.
     {{.Prompt}} {{.HelpName}} play resource

  5. Generate prometheus config for cluster metrics.
     {{.Prompt}} {{.HelpName}} play cluster
`,
}

// PrometheusConfig - container to hold the top level scrape config.
type PrometheusConfig struct {
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs,omitempty"`
}

// String colorized prometheus config yaml.
func (c PrometheusConfig) String() string {
	b, e := yaml.Marshal(c)
	fatalIf(probe.NewError(e), "Unable to generate Prometheus config")

	return console.Colorize("yaml", string(b))
}

// JSON jsonified prometheus config.
func (c PrometheusConfig) JSON() string {
	jsonMessageBytes, e := json.MarshalIndent(c.ScrapeConfigs[0], "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(jsonMessageBytes)
}

// StatConfig - container to hold the targets config.
type StatConfig struct {
	Targets []string `yaml:",flow" json:"targets"`
}

// String colorized stat config yaml.
func (t StatConfig) String() string {
	b, e := yaml.Marshal(t)
	fatalIf(probe.NewError(e), "Unable to generate Prometheus config")

	return console.Colorize("yaml", string(b))
}

// JSON jsonified stat config.
func (t StatConfig) JSON() string {
	jsonMessageBytes, e := json.MarshalIndent(t.Targets, "", " ")
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(jsonMessageBytes)
}

// ScrapeConfig configures a scraping unit for Prometheus.
type ScrapeConfig struct {
	JobName       string       `yaml:"job_name" json:"jobName"`
	BearerToken   string       `yaml:"bearer_token,omitempty" json:"bearerToken,omitempty"`
	MetricsPath   string       `yaml:"metrics_path,omitempty" json:"metricsPath"`
	Scheme        string       `yaml:"scheme,omitempty" json:"scheme"`
	StaticConfigs []StatConfig `yaml:"static_configs,omitempty" json:"staticConfigs"`
}

const (
	defaultPrometheusJWTExpiry = 100 * 365 * 24 * time.Hour
)

// checkAdminPrometheusSyntax - validate all the passed arguments
func checkAdminPrometheusSyntax(ctx context.Context, cmd *cli.Command) {
	args := cmd.Args()
	if len(args.Slice()) == 0 || len(args.Slice()) > 2 {
		showCommandHelpAndExit(ctx, cmd, 1) // last argument is exit code
	}
}

func generatePrometheusConfig(ctx context.Context, cmd *cli.Command) error {
	// Get the alias parameter from cli
	args := cmd.Args()
	alias := cleanAlias(args.Get(0))

	if !isValidAlias(alias) {
		fatalIf(errInvalidAlias(alias), "Invalid alias.")
	}

	hostConfig := mustGetHostConfig(alias)
	if hostConfig == nil {
		fatalIf(errInvalidAliasedURL(alias), "No such alias `"+alias+"` found.")
		return nil
	}

	u, e := url.Parse(hostConfig.URL)
	if e != nil {
		return e
	}

	metricsSubSystem := args.Get(1)
	apiVer := cmd.String("api-version")
	jobName := defaultJobName
	metricsPath := ""

	switch apiVer {
	case "v2":
		if metricsSubSystem == "" {
			metricsSubSystem = "cluster"
		}
		validateV2Args(ctx, cmd, metricsSubSystem)
		if metricsSubSystem != "cluster" {
			jobName = defaultJobName + "-" + metricsSubSystem
		}
		metricsPath = metricsV2BasePath + "/" + metricsSubSystem
	case "v3":
		bucket := cmd.String("bucket")
		validateV3Args(metricsSubSystem, bucket)
		metricsPath = getMetricsV3Path(metricsSubSystem, bucket)
		if metricsSubSystem != "" {
			jobName = defaultJobName + "-" + metricsSubSystem
		}
	default:
		fatalIf(errInvalidArgument().Trace(), "Invalid api version `"+apiVer+"`")
	}

	config := PrometheusConfig{
		ScrapeConfigs: []ScrapeConfig{
			{
				JobName:     jobName,
				MetricsPath: metricsPath,
				StaticConfigs: []StatConfig{
					{
						Targets: []string{""},
					},
				},
			},
		},
	}

	if !cmd.Bool("public") {
		token, e := getPrometheusToken(hostConfig)
		if e != nil {
			return e
		}
		// Setting the values
		config.ScrapeConfigs[0].BearerToken = token
	}
	config.ScrapeConfigs[0].Scheme = u.Scheme
	config.ScrapeConfigs[0].StaticConfigs[0].Targets[0] = u.Host

	printMsg(config)

	return nil
}

// mainAdminPrometheus is the handle for "mc admin prometheus generate" sub-command.
func mainAdminPrometheusGenerate(ctx context.Context, cmd *cli.Command) error {
	console.SetColor("yaml", color.New(color.FgGreen))

	checkAdminPrometheusSyntax(ctx, cmd)

	if err := generatePrometheusConfig(ctx, cmd); err != nil {
		return nil
	}

	return nil
}
