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

func TestRemoteTemplateResolve(t *testing.T) {
	tmpl, err := ParseRemoteTemplate("git@github.com:{repository}.git")
	if err != nil {
		t.Fatal(err)
	}
	got, err := tmpl.Resolve("arcoris/foundation", "foundation")
	if err != nil {
		t.Fatal(err)
	}
	if got != "git@github.com:arcoris/foundation.git" {
		t.Fatalf("Resolve() = %q", got)
	}

	tmpl, err = ParseRemoteTemplate("file:///tmp/{owner}/{name}.git")
	if err != nil {
		t.Fatal(err)
	}
	got, err = tmpl.Resolve("arcoris/control", "control")
	if err != nil {
		t.Fatal(err)
	}
	if got != "file:///tmp/arcoris/control.git" {
		t.Fatalf("Resolve(owner/name) = %q", got)
	}
}

func TestRemoteTemplateRejectsInvalidPlaceholders(t *testing.T) {
	tests := []string{
		"",
		" {repository}",
		"https://example/{missing}.git",
		"https://example/{repository.git",
		"https://example/repository}.git",
		"https://example/\n{repository}.git",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if _, err := ParseRemoteTemplate(tt); err == nil {
				t.Fatalf("ParseRemoteTemplate(%q) error = nil", tt)
			}
		})
	}
}
