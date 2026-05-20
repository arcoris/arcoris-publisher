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
	"encoding/json"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// Env runs go env -json and decodes the reported environment map.
func (t *Toolchain) Env(ctx context.Context, opts goport.EnvOptions) (goport.EnvResult, error) {
	common := goport.CommonOptions{GoBinary: opts.GoBinary, Env: opts.Env, Timeout: opts.Timeout}
	result, err := t.runner.Run(ctx, t.command("", []string{"env", "-json"}, common))
	if err != nil {
		return goport.EnvResult{}, wrapGoError(goport.CodeCommandFailed, "go env failed", result, err)
	}
	values := map[string]string{}
	if err := json.Unmarshal(result.Stdout, &values); err != nil {
		return goport.EnvResult{}, goError(goport.CodeCommandFailed, "go env output could not be parsed", err, nil)
	}
	return goport.EnvResult{Values: values}, nil
}
