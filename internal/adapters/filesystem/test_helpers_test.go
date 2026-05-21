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

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

// writeFile creates parent directories and writes a small fixture file.
func writeFile(t *testing.T, path string, data string) {
	t.Helper()
	must(t, os.MkdirAll(filepath.Dir(path), 0o755))
	must(t, os.WriteFile(path, []byte(data), 0o644))
}

// must fails the current test when err is non-nil.
func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// assertPortCode verifies that err is a structured port error with code.
func assertPortCode(t *testing.T, err error, code porterr.Code) {
	t.Helper()
	var perr *porterr.Error
	if !errors.As(err, &perr) {
		t.Fatalf("expected port error %s, got %T %v", code, err, err)
	}
	if perr.Code != code {
		t.Fatalf("expected code %s, got %s", code, perr.Code)
	}
}

// assertExists verifies that a fixture path is present after an operation.
func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

// assertMissing verifies that a fixture path is absent after an operation.
func assertMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be missing, stat err = %v", path, err)
	}
}
