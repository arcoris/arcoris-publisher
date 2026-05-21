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

import "testing"

func TestDetectFormat(t *testing.T) {
	cases := map[string]Format{
		"arcpub.yaml": FormatYAML,
		"arcpub.yml":  FormatYAML,
		"arcpub.json": FormatJSON,
		"ARCPUB.YAML": FormatYAML,
		"module.JSON": FormatJSON,
	}
	for path, want := range cases {
		got, err := DetectFormat(path)
		if err != nil {
			t.Fatalf("DetectFormat(%q): %v", path, err)
		}
		if got != want {
			t.Fatalf("DetectFormat(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestDetectFormatRejectsUnsupportedExtension(t *testing.T) {
	if _, err := DetectFormat("arcpub.toml"); err == nil {
		t.Fatal("expected unsupported format error")
	}
}

func TestStagingManifestNames(t *testing.T) {
	names := StagingManifestNames()
	if len(names) == 0 || names[0] != "arcpub.yaml" {
		t.Fatalf("unexpected names: %#v", names)
	}
}

func TestModuleManifestNames(t *testing.T) {
	names := ModuleManifestNames()
	if len(names) == 0 || names[0] != "arcpub.module.yaml" {
		t.Fatalf("unexpected names: %#v", names)
	}
}
