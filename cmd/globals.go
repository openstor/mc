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

// Package cmd contains all the global variables and constants. ONLY TO BE ACCESSED VIA GET/SET FUNCTIONS.
package cmd

import (
	"context"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/muesli/termenv"
	"github.com/openstor/madmin-go/v4"
	"github.com/openstor/pkg/v3/console"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/http/httpguts"
)

const (
	globalMCConfigVersion = "10"

	globalMCConfigFile = "config.json"
	globalMCCertsDir   = "certs"
	globalMCCAsDir     = "CAs"

	// session config and shared urls related constants
	globalSessionDir           = "session"
	globalSharedURLsDataDir    = "share"
	globalSessionConfigVersion = "8"

	// Profile directory for dumping profiler outputs.
	globalProfileDir = "profile"

	// Global error exit status.
	globalErrorExitStatus = 1

	// Global CTRL-C (SIGINT, #2) exit status.
	globalCancelExitStatus = 130

	// Global SIGKILL (#9) exit status.
	globalKillExitStatus = 137

	// Global SIGTERM (#15) exit status
	globalTerminatExitStatus = 143
)

var (
	globalQuiet        = false               // Quiet flag set via command line
	globalJSON         = false               // Json flag set via command line
	globalJSONLine     = false               // Print json as single line.
	globalDebug        = false               // Debug flag set via command line
	globalNoColor      = false               // No Color flag set via command line
	globalInsecure     = false               // Insecure flag set via command line
	globalResolvers    map[string]netip.Addr // Custom mappings from HOST[:PORT] to IP
	globalAirgapped    = false               // Airgapped flag set via command line
	globalSubnetConfig []madmin.SubsysConfig // Subnet config

	// GlobalDevMode is set to true if the program is running in development mode
	GlobalDevMode = false

	// GlobalTrapSignals is set to true if need to trap the registered signals and cancel the global context.
	GlobalTrapSignals = true

	// GlobalSubnetProxyURL is the proxy to be used for communication with subnet
	GlobalSubnetProxyURL *url.URL

	globalConnReadDeadline  time.Duration
	globalConnWriteDeadline time.Duration

	globalLimitUpload   uint64
	globalLimitDownload uint64

	globalContext, globalCancel = context.WithCancel(context.Background())

	globalCustomHeader http.Header
)

var (
	// Terminal height/width, zero if not found
	globalTermWidth, globalTermHeight int

	globalDisablePagerEnv       = "DISABLE_PAGER"
	globalDisablePagerFlag      = "--disable-pager"
	globalDisablePagerFlagShort = "--dp"
	globalPagerDisabled         = false
	globalHelpPager             *termPager

	// CA root certificates, a nil value means system certs pool will be used
	globalRootCAs *x509.CertPool
)

func parsePagerDisableFlag(args []string) {
	globalPagerDisabled, _ = strconv.ParseBool(os.Getenv(envPrefix + globalDisablePagerEnv))
	for _, arg := range args {
		if arg == globalDisablePagerFlag || arg == globalDisablePagerFlagShort {
			globalPagerDisabled = true
		}
	}
}

// Set global states. NOTE: It is deliberately kept monolithic to ensure we dont miss out any flags.
func setGlobalsFromContext(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	quiet := cmd.Bool("quiet")       // || cmd.GlobalBool("quiet")
	debug := cmd.Bool("debug")       // || cmd.GlobalBool("debug")
	json := cmd.Bool("json")         // || cmd.GlobalBool("json")
	noColor := cmd.Bool("no-color")  // || cmd.GlobalBool("no-color")
	insecure := cmd.Bool("insecure") // || cmd.GlobalBool("insecure")
	devMode := cmd.Bool("dev")       // || cmd.GlobalBool("dev")
	airgapped := cmd.Bool("airgap")  // || cmd.GlobalBool("airgap")

	globalQuiet = globalQuiet || quiet
	globalDebug = globalDebug || debug
	globalJSONLine = !isTerminal() && json
	globalJSON = globalJSON || json
	globalNoColor = globalNoColor || noColor || globalJSONLine
	globalInsecure = globalInsecure || insecure
	GlobalDevMode = GlobalDevMode || devMode
	globalAirgapped = globalAirgapped || airgapped

	// Disable colorified messages if requested.
	if globalNoColor || globalQuiet {
		console.SetColorOff()
		lipgloss.SetColorProfile(termenv.Ascii)
	}

	globalConnReadDeadline = cmd.Duration("conn-read-deadline")
	// if globalConnReadDeadline <= 0 {
	// 	globalConnReadDeadline = cmd.GlobalDuration("conn-read-deadline")
	// }

	globalConnWriteDeadline = cmd.Duration("conn-write-deadline")
	// if globalConnWriteDeadline <= 0 {
	// 	globalConnWriteDeadline = cmd.GlobalDuration("conn-write-deadline")
	// }

	limitUploadStr := cmd.String("limit-upload")
	// if limitUploadStr == "" {
	// 	limitUploadStr = cmd.GlobalString("limit-upload")
	// }
	if limitUploadStr != "" {
		var e error
		globalLimitUpload, e = humanize.ParseBytes(limitUploadStr)
		if e != nil {
			return ctx, e
		}
	}

	limitDownloadStr := cmd.String("limit-download")
	// if limitDownloadStr == "" {
	// 	limitDownloadStr = cmd.GlobalString("limit-download")
	// }

	if limitDownloadStr != "" {
		var e error
		globalLimitDownload, e = humanize.ParseBytes(limitDownloadStr)
		if e != nil {
			return ctx, e
		}
	}

	dnsEntries := cmd.StringSlice("resolve")
	if len(dnsEntries) > 0 {
		globalResolvers = make(map[string]netip.Addr, len(dnsEntries))

		// Each entry is a HOST[:PORT]=IP pair. This is very similar to cURL's syntax.
		for _, e := range dnsEntries {
			i := strings.IndexByte(e, '=')
			if i < 0 {
				return ctx, fmt.Errorf("invalid DNS resolve entry %s", e)
			}

			if strings.ContainsRune(e[:i], ':') {
				if _, _, err := net.SplitHostPort(e[:i]); err != nil {
					return ctx, fmt.Errorf("invalid DNS resolve entry %s: %v", e, err)
				}
			}

			host := e[:i]
			addr, err := netip.ParseAddr(e[i+1:])
			if err != nil {
				return ctx, fmt.Errorf("invalid DNS resolve entry %s: %v", e, err)
			}
			globalResolvers[host] = addr
		}
	}

	customHeaders := cmd.StringSlice("custom-header")
	if len(customHeaders) > 0 {
		globalCustomHeader = make(http.Header)
		for _, header := range customHeaders {
			i := strings.IndexByte(header, ':')
			if i <= 0 {
				return ctx, fmt.Errorf("invalid custom header entry %s", header)
			}
			h := strings.TrimSpace(header[:i])
			hv := strings.TrimSpace(header[i+1:])
			if !httpguts.ValidHeaderFieldName(h) || !httpguts.ValidHeaderFieldValue(hv) {
				return ctx, fmt.Errorf("invalid custom header entry %s", header)
			}
			globalCustomHeader.Add(h, hv)
		}
	}

	return ctx, nil
}
