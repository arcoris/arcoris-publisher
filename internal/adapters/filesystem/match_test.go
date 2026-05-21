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

func TestShouldInclude(t *testing.T) {
	tests := []struct {
		name    string
		rel     string
		include []string
		exclude []string
		want    bool
	}{
		{name: "default", rel: "pkg/x.go", want: true},
		{name: "git dir", rel: ".git/config", want: false},
		{name: "include match", rel: "pkg/x.go", include: []string{"**/*.go"}, want: true},
		{name: "include miss", rel: "pkg/x.txt", include: []string{"**/*.go"}, want: false},
		{name: "exclude wins", rel: "pkg/x.go", include: []string{"**/*.go"}, exclude: []string{"pkg/**"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldInclude(tt.rel, tt.include, tt.exclude); got != tt.want {
				t.Fatalf("shouldInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlash(t *testing.T) {
	if got := slash(`.\pkg\x.go`); got != "pkg/x.go" {
		t.Fatalf("slash() = %q, want pkg/x.go", got)
	}
}

func TestMatchPatternDialects(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{name: "empty", pattern: "", path: "pkg/x.go", want: false},
		{name: "exact", pattern: "pkg/x.go", path: "pkg/x.go", want: true},
		{name: "subtree", pattern: "pkg/**", path: "pkg/nested/x.go", want: true},
		{name: "basename anywhere", pattern: "**/*.go", path: "pkg/x.go", want: true},
		{name: "basename glob", pattern: "*.go", path: "pkg/x.go", want: true},
		{name: "miss", pattern: "*.txt", path: "pkg/x.go", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchPattern(tt.pattern, tt.path); got != tt.want {
				t.Fatalf("matchPattern(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
			}
		})
	}
}
