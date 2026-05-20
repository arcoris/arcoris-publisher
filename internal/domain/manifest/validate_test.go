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

func TestParseModuleName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "lower", value: "foundation"},
		{name: "dash", value: "control-plane"},
		{name: "underscore", value: "control_plane"},
		{name: "empty", value: "", wantErr: true},
		{name: "upper", value: "Foundation", wantErr: true},
		{name: "slash", value: "foundation/core", wantErr: true},
		{name: "starts dash", value: "-foundation", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseModuleName(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseModuleName(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestParseSourceDir(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    SourceDir
		wantErr bool
	}{
		{name: "normal", value: "staging/src/arcoris.dev/foundation", want: SourceDir("staging/src/arcoris.dev/foundation")},
		{name: "clean", value: "staging/./src/arcoris.dev/foundation", want: SourceDir("staging/src/arcoris.dev/foundation")},
		{name: "empty", value: "", wantErr: true},
		{name: "absolute", value: "/tmp/module", wantErr: true},
		{name: "escape", value: "../module", wantErr: true},
		{name: "backslash", value: "staging\\module", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSourceDir(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseSourceDir(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("ParseSourceDir(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestParseRepositoryRef(t *testing.T) {
	ref, err := ParseRepositoryRef("arcoris/foundation")
	if err != nil {
		t.Fatalf("ParseRepositoryRef() error = %v", err)
	}
	if ref.Owner() != "arcoris" || ref.Name() != "foundation" {
		t.Fatalf("Owner/Name = %q/%q", ref.Owner(), ref.Name())
	}
	invalid := []string{"", "arcoris", "arcoris/foundation/extra", "../foundation", "arcoris/..", "ar coris/foundation"}
	for _, value := range invalid {
		if _, err := ParseRepositoryRef(value); err == nil {
			t.Fatalf("ParseRepositoryRef(%q) error = nil, want error", value)
		}
	}
}

func TestValidationErrorError(t *testing.T) {
	err := (&ValidationError{Issues: []Issue{{Path: "version", Message: "bad"}}}).Error()
	if err != "version: bad" {
		t.Fatalf("Error() = %q", err)
	}
	err = (&ValidationError{Issues: []Issue{{Message: "one"}, {Message: "two"}}}).Error()
	if err != "manifest validation failed with 2 issues" {
		t.Fatalf("Error() = %q", err)
	}
}

func TestParseModulePath(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "domain", value: "arcoris.dev/foundation"},
		{name: "dotted", value: "arcoris.dev"},
		{name: "empty", value: "", wantErr: true},
		{name: "space", value: "arcoris.dev/foundation module", wantErr: true},
		{name: "backslash", value: "arcoris.dev\\foundation", wantErr: true},
		{name: "plain", value: "foundation", wantErr: true},
		{name: "absolute", value: "/arcoris.dev/foundation", wantErr: true},
		{name: "traversal", value: "../foundation", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseModulePath(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseModulePath(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestZeroRepositoryRefOwnerName(t *testing.T) {
	var ref RepositoryRef
	if ref.Owner() != "" || ref.Name() != "" {
		t.Fatalf("zero ref Owner/Name = %q/%q", ref.Owner(), ref.Name())
	}
}

func TestValidationErrorEmpty(t *testing.T) {
	if got := (*ValidationError)(nil).Error(); got != "manifest validation failed" {
		t.Fatalf("nil ValidationError Error() = %q", got)
	}
	if got := (&ValidationError{}).Error(); got != "manifest validation failed" {
		t.Fatalf("empty ValidationError Error() = %q", got)
	}
}
