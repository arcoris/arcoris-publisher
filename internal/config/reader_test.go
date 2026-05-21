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

package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestOSReaderReadsAndChecksExistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "arcpub.yaml")
	if err := os.WriteFile(path, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	reader := OSReader{}
	exists, err := reader.Exists(context.Background(), path)
	if err != nil || !exists {
		t.Fatalf("Exists = %v, %v", exists, err)
	}
	data, err := reader.ReadFile(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q", data)
	}
}

func TestOSReaderHonorsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (OSReader{}).Exists(ctx, "anything"); err == nil {
		t.Fatal("expected cancellation")
	}
}

func TestOSReaderReadFileHonorsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (OSReader{}).ReadFile(ctx, "anything"); err == nil {
		t.Fatal("expected cancellation")
	}
}

func TestOSReaderExistsReturnsFalseForMissingPath(t *testing.T) {
	exists, err := (OSReader{}).Exists(
		context.Background(),
		filepath.Join(t.TempDir(), "missing"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("expected missing path")
	}
}
