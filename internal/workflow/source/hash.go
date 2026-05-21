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

package source

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Hash is a deterministic sha256 content identifier produced by source
// inspection.
type Hash string

// String returns the hash string.
func (h Hash) String() string { return string(h) }

// IsZero reports whether no hash was computed.
func (h Hash) IsZero() bool { return h == "" }

// hashBytes hashes typed ordered parts with NUL separators to avoid accidental
// concatenation collisions between adjacent fields.
func hashBytes(kind string, parts ...string) Hash {
	h := sha256.New()
	_, _ = h.Write([]byte(kind))
	_, _ = h.Write([]byte{0})
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
		_, _ = h.Write([]byte{0})
	}
	return Hash("sha256:" + hex.EncodeToString(h.Sum(nil)))
}

// combineHashes combines non-zero child hashes into one parent hash.
func combineHashes(kind string, values []Hash) Hash {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		if value.IsZero() {
			continue
		}
		parts = append(parts, value.String())
	}
	if len(parts) == 0 {
		return ""
	}
	return hashBytes(kind, strings.Join(parts, "\n"))
}
