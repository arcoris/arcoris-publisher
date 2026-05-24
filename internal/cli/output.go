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
	"fmt"
	"io"

	"arcoris.dev/arcoris-publisher/internal/report"
)

func newRenderer(opts report.Options) report.Renderer {
	return report.New(opts)
}

func writeCLIError(w io.Writer, err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintf(w, "arcpub: %v\n", err)
}

func writeUsage(w io.Writer) {
	_, _ = fmt.Fprint(w, `ARCORIS Publisher

Usage:
  arcpub plan    --manifest arcpub.yaml --version v0.3.0 [--output text|json]
  arcpub verify  --manifest arcpub.yaml --version v0.3.0 [workflow flags]
  arcpub publish --manifest arcpub.yaml --version v0.3.0 [workflow flags]
  arcpub version [--output text|json]
  arcpub help

Workflow flags:
  --source-repo PATH       Source Git checkout root. Default: .
  --staging-dir PATH       Staging directory containing module sources. Default: .
  --target-root PATH       Target worktree root. Default: .arcpub/targets
  --dry-run                Verify but do not publish refs.

Output flags:
  --output text|json       Report format. Default: text.
  --include-local-paths    Include local absolute filesystem paths in reports.
  --pretty                 Pretty JSON when --output=json.
  --compact                Compact JSON when --output=json.
`)
}
