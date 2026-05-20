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

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// List runs go list and optionally parses newline-delimited JSON package data.
func (t *Toolchain) List(ctx context.Context, moduleDir string, opts goport.ListOptions) (goport.ListResult, error) {
	result, err := t.runner.Run(ctx, t.command(moduleDir, listArgs(opts), opts.CommonOptions))
	out, parseErr := parseListResult(result.Stdout, result.Stderr, opts)
	if parseErr != nil {
		return out, parseErr
	}
	if err != nil {
		return out, wrapGoError(goport.CodeListFailed, "go list failed", result, err)
	}
	return out, nil
}
