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
	"path/filepath"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestRemoveAllRemovesPathInsideSafetyRoot(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "dir")
	writeFile(t, filepath.Join(dir, "file.txt"), "content")

	if err := New().RemoveAll(context.Background(), dir, fsport.RemoveOptions{SafetyRoot: root}); err != nil {
		t.Fatalf("RemoveAll() error = %v", err)
	}
	assertMissing(t, dir)
}

func TestRemoveAllAllowMissingAndUnsafe(t *testing.T) {
	root := t.TempDir()
	fs := New()

	err := fs.RemoveAll(context.Background(), filepath.Join(root, "missing"), fsport.RemoveOptions{SafetyRoot: root, AllowMissing: true})
	if err != nil {
		t.Fatalf("RemoveAll() missing allowed error = %v", err)
	}
	err = fs.RemoveAll(context.Background(), filepath.Dir(root), fsport.RemoveOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodeUnsafeRemove)
}

func TestValidateRemoveTargetMissingRequiresAllowMissing(t *testing.T) {
	err := validateRemoveTarget(filepath.Join(t.TempDir(), "missing"), fsport.RemoveOptions{})
	assertPortCode(t, err, fsport.CodePathNotFound)
}

func TestRemoveAllMissingRequiresAllowMissing(t *testing.T) {
	root := t.TempDir()
	err := New().RemoveAll(context.Background(), filepath.Join(root, "missing"), fsport.RemoveOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodePathNotFound)
}

func TestRemoveAllContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := New().RemoveAll(ctx, filepath.Join(t.TempDir(), "dir"), fsport.RemoveOptions{}); err == nil {
		t.Fatalf("RemoveAll() should return context cancellation")
	}
}
