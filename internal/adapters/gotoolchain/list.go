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
	"strings"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// List runs go list and optionally parses newline-delimited JSON package data.
func (t *Toolchain) List(ctx context.Context, moduleDir string, opts goport.ListOptions) (goport.ListResult, error) {
	patterns := opts.Patterns
	if len(patterns) == 0 {
		patterns = []string{"./..."}
	}
	args := []string{"list"}
	if opts.JSON {
		args = append(args, "-json")
	}
	if opts.Deps {
		args = append(args, "-deps")
	}
	if opts.Test {
		args = append(args, "-test")
	}
	if len(opts.Tags) > 0 {
		args = append(args, "-tags", strings.Join(opts.Tags, ","))
	}
	args = append(args, patterns...)
	result, err := t.runner.Run(ctx, t.command(moduleDir, args, opts.CommonOptions))
	out := goport.ListResult{Stdout: result.Stdout, Stderr: result.Stderr}
	if opts.JSON && len(result.Stdout) > 0 {
		packages, parseErr := parsePackages(result.Stdout)
		if parseErr != nil {
			return out, goError(goport.CodeListFailed, "go list output could not be parsed", parseErr, nil)
		}
		out.Packages = packages
	}
	if err != nil {
		return out, wrapGoError(goport.CodeListFailed, "go list failed", result, err)
	}
	return out, nil
}
