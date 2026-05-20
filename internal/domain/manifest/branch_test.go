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

package manifest

import "testing"

func TestParseBranchName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "main", value: "main"},
		{name: "release", value: "release/v1"},
		{name: "empty", value: "", wantErr: true},
		{name: "space", value: "feature branch", wantErr: true},
		{name: "dash", value: "-main", wantErr: true},
		{name: "traversal", value: "../main", wantErr: true},
		{name: "invalid char", value: "main~old", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseBranchName(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseBranchName(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestNewBranchMapping(t *testing.T) {
	mapping, err := NewBranchMapping(BranchMappingSpec{Source: "main", Target: "release/main"})
	if err != nil {
		t.Fatalf("NewBranchMapping() error = %v", err)
	}
	if mapping.Source() != BranchName("main") || mapping.Target() != BranchName("release/main") {
		t.Fatalf("mapping = %q -> %q", mapping.Source(), mapping.Target())
	}
	if got := mapping.Spec(); got.Source != "main" || got.Target != "release/main" {
		t.Fatalf("Spec() = %#v", got)
	}
}

func TestNewBranchMappingRejectsInvalidTarget(t *testing.T) {
	if _, err := NewBranchMapping(BranchMappingSpec{Source: "main", Target: "target branch"}); err == nil {
		t.Fatalf("NewBranchMapping() error = nil, want error")
	}
}
