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
	"strconv"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// testArgs renders typed TestOptions into a go test command line.
func testArgs(opts goport.TestOptions) []string {
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
	args = appendBuildTags(args, opts.Tags)
	return append(args, defaultPatterns(opts.Patterns)...)
}
