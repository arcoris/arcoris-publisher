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
	"errors"
	"os"
	"path/filepath"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestPathSafetyHelpers(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "child")
	outside := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-other")

	if err := ensureInside(outside, ""); err != nil {
		t.Fatalf("ensureInside(empty root) error = %v", err)
	}
	if err := ensureInside(root, root); err != nil {
		t.Fatalf("ensureInside(root) error = %v", err)
	}
	if err := ensureInside(inside, root); err != nil {
		t.Fatalf("ensureInside(inside) error = %v", err)
	}
	if err := ensureInside(outside, root); err == nil {
		t.Fatalf("ensureInside(outside) should fail")
	}
	if _, err := absClean(" "); err == nil {
		t.Fatalf("absClean(empty) should fail")
	}
	if err := ensureInside(inside, " "); err == nil {
		t.Fatalf("ensureInside(empty root text) should fail")
	}
	if err := ensureSafeRemove(string(filepath.Separator), ""); err == nil {
		t.Fatalf("ensureSafeRemove(root) should fail")
	}
}

func TestPathPolicyHelpers(t *testing.T) {
	wrappedNotExist := &os.PathError{Op: "stat", Path: "missing", Err: os.ErrNotExist}
	if !isNotExist(os.ErrNotExist) || !isNotExist(wrappedNotExist) {
		t.Fatalf("isNotExist() should recognize wrapped not-exist errors")
	}
	if isNotExist(errors.New("other")) {
		t.Fatalf("isNotExist() should reject unrelated errors")
	}
	if got := symlinkMode(""); got != fsport.SymlinkReject {
		t.Fatalf("symlinkMode(empty) = %s, want reject", got)
	}
	if got := symlinkMode(fsport.SymlinkPreserve); got != fsport.SymlinkPreserve {
		t.Fatalf("symlinkMode(preserve) = %s, want preserve", got)
	}
}
