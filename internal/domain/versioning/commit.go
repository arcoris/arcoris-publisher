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

package versioning

import (
	"fmt"
	"strings"
	"unicode"
)

// normalizeCommit validates and lowercases a source commit hash.
func normalizeCommit(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("commit is required")
	}
	if len(value) < minimumSnapshotCommitLength {
		return "", fmt.Errorf("commit must contain at least %d hexadecimal characters", minimumSnapshotCommitLength)
	}
	for _, r := range value {
		if !isHex(r) {
			return "", fmt.Errorf("commit contains non-hexadecimal character %q", r)
		}
	}
	return strings.ToLower(value), nil
}

// isHex reports whether r is an ASCII hexadecimal digit.
func isHex(r rune) bool {
	return unicode.IsDigit(r) || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}
