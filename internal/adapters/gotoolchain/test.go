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

package gotoolchain

import (
	"context"
	"strconv"
	"strings"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// Test runs go test with typed test options.
func (t *Toolchain) Test(ctx context.Context, moduleDir string, opts goport.TestOptions) (goport.TestResult, error) {
	patterns := opts.Patterns
	if len(patterns) == 0 {
		patterns = []string{"./..."}
	}
	args := []string{"test"}
	if opts.Race {
		args = append(args, "-race")
	}
	if opts.Count > 0 {
		args = append(args, "-count", strconv.Itoa(opts.Count))
	}
	if opts.Short {
		args = append(args, "-short")
	}
	if opts.Run != "" {
		args = append(args, "-run", opts.Run)
	}
	if opts.Verbose {
		args = append(args, "-v")
	}
	if len(opts.Tags) > 0 {
		args = append(args, "-tags", strings.Join(opts.Tags, ","))
	}
	args = append(args, patterns...)
	result, err := t.runner.Run(ctx, t.command(moduleDir, args, opts.CommonOptions))
	out := goport.TestResult{Stdout: result.Stdout, Stderr: result.Stderr}
	if err != nil {
		return out, wrapGoError(goport.CodeTestFailed, "go test failed", result, err)
	}
	return out, nil
}
