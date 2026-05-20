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

import "strings"

// defaultPatterns applies ./... when an operation accepts package patterns and
// the caller did not provide an explicit set.
func defaultPatterns(patterns []string) []string {
	if len(patterns) == 0 {
		return []string{"./..."}
	}
	return append([]string(nil), patterns...)
}

// appendBuildTags appends one comma-joined -tags value when tags are present.
//
// The Go command expects a single -tags flag whose value is a comma-separated
// list. Repeating -tags can produce surprising behavior, so both list and test
// share this helper.
func appendBuildTags(args []string, tags []string) []string {
	if len(tags) == 0 {
		return args
	}
	return append(args, "-tags", strings.Join(tags, ","))
}
