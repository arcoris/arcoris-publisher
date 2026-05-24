// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"arcoris.dev/arcoris-publisher/internal/buildinfo"
	"arcoris.dev/arcoris-publisher/internal/report"
	"github.com/spf13/cobra"
)

// BuildInfoFunc returns normalized publisher build metadata.
type BuildInfoFunc func() buildinfo.Info

type versionReport struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Dirty   string `json:"dirty"`
}

func (c CLI) newVersionCommand(output *outputFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print publisher build metadata",
		Args:  noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.executeVersion(outputForCommand(cmd, *output), cmd.OutOrStdout())
		},
	}
}

func (c CLI) executeVersion(output outputFlags, stdout io.Writer) error {
	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}

	info := c.deps.BuildInfo()
	value := versionReport{
		Version: info.Version(),
		Commit:  info.Commit(),
		Date:    info.Date(),
		Dirty:   info.Dirty(),
	}

	switch reportOptions.Format {
	case report.FormatText:
		_, err = fmt.Fprintf(
			stdout,
			"arcpub %s\n  commit: %s\n  date:   %s\n  dirty:  %s\n",
			value.Version,
			value.Commit,
			value.Date,
			value.Dirty,
		)
	case report.FormatJSON:
		var data []byte
		if reportOptions.Pretty {
			data, err = json.MarshalIndent(value, "", "  ")
		} else {
			data, err = json.Marshal(value)
		}
		if err == nil {
			_, err = stdout.Write(append(data, '\n'))
		}
	default:
		err = &Error{Code: CodeInvalidFlags, Message: "unsupported version output format"}
	}
	if err != nil {
		return &Error{Code: CodeReportFailed, Message: "render version failed", Cause: err}
	}
	return nil
}
