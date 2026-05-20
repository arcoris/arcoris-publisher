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
	"context"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// TreeHash walks root and computes a deterministic content digest.
//
// The walker gathers normalized entries first and sorts them before hashing so
// filesystem iteration order cannot affect the digest. File content is hashed
// separately, then folded into the tree hash with relative path and selected
// metadata.
func (fs *FileSystem) TreeHash(ctx context.Context, root string, opts fsport.TreeHashOptions) (fsport.TreeHash, error) {
	entries, err := collectHashEntries(ctx, root, opts)
	if err != nil {
		if isPortError(err) {
			return "", err
		}
		if isNotExist(err) {
			return "", pathError(fsport.CodePathNotFound, "tree hash root not found", err, root)
		}
		return "", fsError(fsport.CodeTreeHashFailed, "tree hash failed", err, porterrDetails("root", root))
	}
	return renderTreeHash(entries, opts.IncludeFileMode), nil
}
