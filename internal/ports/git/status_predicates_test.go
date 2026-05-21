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

package git

import "testing"

func TestStatusHasEntries(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{name: "empty", status: Status{Clean: true}, want: false},
		{name: "entry", status: Status{Entries: []StatusEntry{{Path: "go.mod", Code: "M "}}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.HasEntries(); got != tt.want {
				t.Fatalf("HasEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusIsDirty(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{name: "explicit clean", status: Status{Clean: true}, want: false},
		{name: "clean flag false", status: Status{}, want: true},
		{name: "entries override clean flag", status: Status{Clean: true, Entries: []StatusEntry{{Path: "go.mod"}}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsDirty(); got != tt.want {
				t.Fatalf("IsDirty() = %v, want %v", got, tt.want)
			}
		})
	}
}
