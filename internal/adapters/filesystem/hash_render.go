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

package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// renderTreeHash sorts entries and folds them into the public TreeHash value.
func renderTreeHash(entries []hashEntry, includeMode bool) fsport.TreeHash {
	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })
	sum := sha256.New()
	for _, entry := range entries {
		writeHashEntry(sum, entry, includeMode)
	}
	return fsport.TreeHash("sha256:" + hex.EncodeToString(sum.Sum(nil)))
}

// writeHashEntry writes one normalized entry to the digest stream.
func writeHashEntry(sum interface{ Write([]byte) (int, error) }, entry hashEntry, includeMode bool) {
	if includeMode {
		fmt.Fprintf(sum, "%s %s %o %s\n", entry.kind, entry.rel, entry.mode, entry.content)
		return
	}
	fmt.Fprintf(sum, "%s %s %s\n", entry.kind, entry.rel, entry.content)
}
