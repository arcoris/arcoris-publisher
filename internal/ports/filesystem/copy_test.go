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

import "testing"

func TestCopyTreeResultCopiedEntries(t *testing.T) {
	result := CopyTreeResult{
		FilesCopied:       3,
		DirectoriesCopied: 2,
		FilesSkipped:      7,
		BytesCopied:       128,
	}

	if got := result.CopiedEntries(); got != 5 {
		t.Fatalf("CopiedEntries() = %d, want 5", got)
	}
}

func TestCopyTreeResultEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   CopyTreeResult
		want bool
	}{
		{name: "zero", in: CopyTreeResult{}, want: true},
		{name: "skipped only", in: CopyTreeResult{FilesSkipped: 1}, want: true},
		{name: "bytes only diagnostic", in: CopyTreeResult{BytesCopied: 1}, want: true},
		{name: "file copied", in: CopyTreeResult{FilesCopied: 1}, want: false},
		{name: "directory copied", in: CopyTreeResult{DirectoriesCopied: 1}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Empty(); got != tt.want {
				t.Fatalf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
