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

// listArgs renders typed ListOptions into a go list command line.
func listArgs(opts goport.ListOptions) []string {
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
	args = appendBuildTags(args, opts.Tags)
	return append(args, defaultPatterns(opts.Patterns)...)
}
