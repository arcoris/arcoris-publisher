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

package gitcli

import (
	"strings"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// parseStatus parses NUL-delimited porcelain-v1 output.
//
// Rename and copy records include an extra path field after the primary entry;
// the parser skips that extra field because the port currently models only one
// path per status entry.
func parseStatus(out []byte) []gitport.StatusEntry {
	if len(out) == 0 {
		return nil
	}
	parts := strings.Split(string(out), "\x00")
	entries := make([]gitport.StatusEntry, 0, len(parts))
	for i := 0; i < len(parts); i++ {
		entry, consumesExtraPath := parseStatusPart(parts[i])
		if entry == nil {
			continue
		}
		entries = append(entries, *entry)
		if consumesExtraPath {
			i++
		}
	}
	return entries
}

// parseStatusPart converts one porcelain record into the port status shape.
func parseStatusPart(part string) (*gitport.StatusEntry, bool) {
	if part == "" {
		return nil, false
	}
	if len(part) < 3 {
		return &gitport.StatusEntry{Path: part}, false
	}
	code := part[:2]
	path := strings.TrimPrefix(part[2:], " ")
	return &gitport.StatusEntry{Code: code, Path: path}, statusCodeHasExtraPath(code)
}

// statusCodeHasExtraPath reports whether porcelain output has a second path.
func statusCodeHasExtraPath(code string) bool {
	return strings.HasPrefix(code, "R") || strings.HasPrefix(code, "C")
}
