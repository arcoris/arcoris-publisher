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

package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"arcoris.dev/arcoris-publisher/internal/workflow/source"
)

const absentEntryHash = "absent"

// Entry describes one explicit projection target used for provenance hashing.
//
// It contains only target repository-relative data. Source absolute paths are
// intentionally omitted because provenance may be committed to target
// repositories.
type Entry struct {
	TargetPath string `json:"targetPath"`
	Hash       string `json:"hash"`
	Present    bool   `json:"present"`
}

// EntriesFromSourceModule converts inspected source entries into stable
// provenance entries.
func EntriesFromSourceModule(module source.ModuleSnapshot) []Entry {
	entries := module.Entries()
	out := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, Entry{
			TargetPath: entry.TargetPath().String(),
			Hash:       projectionHashValue(entry.Present(), entry.Hash().String()),
			Present:    entry.Present(),
		})
	}
	return out
}

// ProjectionHash returns a deterministic sha256 digest for target paths and
// inspected entry hashes.
//
// Present entries use their source content hash. Absent optional entries are
// represented as "<targetPath>=absent" so an omitted optional file changes the
// projection relative to a present empty or unhashed file.
func ProjectionHash(entries []Entry) string {
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, entry.TargetPath+"="+projectionHashValue(entry.Present, entry.Hash))
	}
	sort.Strings(lines)

	hash := sha256.New()
	for _, line := range lines {
		_, _ = hash.Write([]byte(line))
		_, _ = hash.Write([]byte{'\n'})
	}

	return "sha256:" + hex.EncodeToString(hash.Sum(nil))
}

func projectionHashValue(present bool, hash string) string {
	if !present {
		return absentEntryHash
	}
	return hash
}
