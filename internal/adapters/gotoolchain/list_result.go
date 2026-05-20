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

import goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"

// parseListResult copies command output and decodes packages when JSON was requested.
func parseListResult(stdout []byte, stderr []byte, opts goport.ListOptions) (goport.ListResult, error) {
	out := goport.ListResult{Stdout: stdout, Stderr: stderr}
	if !opts.JSON || len(stdout) == 0 {
		return out, nil
	}
	packages, err := parsePackages(stdout)
	if err != nil {
		return out, goError(goport.CodeListFailed, "go list output could not be parsed", err, nil)
	}
	out.Packages = packages
	return out, nil
}
