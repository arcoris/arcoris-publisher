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

package verify

import (
	"context"
	"testing"
)

func TestVerifyRejectsInvalidRequest(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).Verify(context.Background(), Request{})

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodeInvalidRequest {
		t.Fatalf("Code = %q", got.Code)
	}
}

func TestTargetWorktreeCheck(t *testing.T) {
	tests := []struct {
		name   string
		fs     fakeReader
		status Status
	}{
		{
			name:   "missing",
			fs:     fakeReader{},
			status: StatusFailed,
		},
		{
			name: "not directory",
			fs: fakeReader{
				paths: map[string]fakePath{"/target": {exists: true}},
			},
			status: StatusFailed,
		},
		{
			name: "directory",
			fs: fakeReader{
				paths: map[string]fakePath{"/target": {exists: true, dir: true}},
			},
			status: StatusPassed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(Dependencies{FS: tt.fs}, Options{})

			check := service.targetWorktreeCheck(context.Background(), "/target")

			if check.Status() != tt.status {
				t.Fatalf("Status() = %s, want %s", check.Status(), tt.status)
			}
			if check.Path() != "/target" {
				t.Fatalf("Path() = %q", check.Path())
			}
		})
	}
}

type fakePath struct {
	exists bool
	dir    bool
	data   []byte
}

type fakeReader struct {
	paths map[string]fakePath
}

func (fs fakeReader) Exists(_ context.Context, name string) (bool, error) {
	path := fs.paths[name]
	return path.exists, nil
}

func (fs fakeReader) IsDir(_ context.Context, name string) (bool, error) {
	path := fs.paths[name]
	return path.dir, nil
}

func (fs fakeReader) ReadFile(_ context.Context, path string) ([]byte, error) {
	data := fs.paths[path].data
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}
