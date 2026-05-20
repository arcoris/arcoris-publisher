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

// ModTidy runs go mod tidy.
//
// Build tags are intentionally not passed here. The Go command's tidy mode does
// not accept the same -tags flag as go list/test, so tags remain a shared option
// for compatibility but are ignored for this operation.
func (t *Toolchain) ModTidy(ctx context.Context, moduleDir string, opts goport.ModTidyOptions) (goport.ModTidyResult, error) {
	args := []string{"mod", "tidy"}
	if opts.Compat != "" {
		args = append(args, "-compat", opts.Compat)
	}
	result, err := t.runner.Run(ctx, t.command(moduleDir, args, opts.CommonOptions))
	out := goport.ModTidyResult{Stdout: result.Stdout, Stderr: result.Stderr}
	if err != nil {
		return out, wrapGoError(goport.CodeModTidyFailed, "go mod tidy failed", result, err)
	}
	return out, nil
}
