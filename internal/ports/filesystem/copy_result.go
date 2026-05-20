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

// CopyTreeResult reports tree-copy counters useful for diagnostics.
type CopyTreeResult struct {
	// FilesCopied is the number of regular files written to the destination.
	FilesCopied int
	// DirectoriesCopied is the number of directories created or reused.
	DirectoriesCopied int
	// FilesSkipped is the number of source files intentionally ignored.
	FilesSkipped int
	// BytesCopied is the total number of regular-file bytes written.
	BytesCopied int64
}

// CopiedEntries reports how many filesystem entries were copied or created.
func (r CopyTreeResult) CopiedEntries() int {
	return r.FilesCopied + r.DirectoriesCopied
}

// Empty reports whether no filesystem entries were copied.
func (r CopyTreeResult) Empty() bool {
	return r.CopiedEntries() == 0
}
