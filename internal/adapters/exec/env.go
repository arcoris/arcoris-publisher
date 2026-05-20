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

package exec

import "strings"

// MergeEnv overlays override assignments onto base environment assignments.
//
// Entries are expected to be in KEY=VALUE form. Valid override entries replace
// the first base entry with the same key while preserving the original order of
// unrelated base entries. Malformed entries are appended unchanged because the
// process port should not silently discard caller-provided environment data.
func MergeEnv(base, override []string) []string {
	out := append([]string(nil), base...)
	for _, entry := range override {
		key, _, ok := strings.Cut(entry, "=")
		if !ok || key == "" {
			out = append(out, entry)
			continue
		}
		replaced := false
		for i, existing := range out {
			existingKey, _, existingOK := strings.Cut(existing, "=")
			if existingOK && existingKey == key {
				out[i] = entry
				replaced = true
				break
			}
		}
		if !replaced {
			out = append(out, entry)
		}
	}
	return out
}
