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

package gotoolchain

import "testing"

func TestParsePackages(t *testing.T) {
	input := []byte(`{
  "ImportPath":"example.com/a",
  "Module":{"Path":"example.com/a","Replace":{"Path":"../a"}},
  "Imports":["fmt"],
  "TestImports":["testing"],
  "Deps":["errors"]
}
{"ImportPath":"example.com/b"}
`)
	packages, err := parsePackages(input)
	if err != nil {
		t.Fatalf("parsePackages() error = %v", err)
	}
	if len(packages) != 2 {
		t.Fatalf("len = %d, want 2", len(packages))
	}
	if packages[0].Module.Replace == nil || packages[0].Module.Replace.Path != "../a" {
		t.Fatalf("replace not parsed: %#v", packages[0].Module)
	}
	if packages[0].Imports[0] != "fmt" || packages[0].TestImports[0] != "testing" || packages[0].Deps[0] != "errors" {
		t.Fatalf("imports not parsed: %#v", packages[0])
	}
}

func TestParsePackagesMalformedJSON(t *testing.T) {
	if _, err := parsePackages([]byte(`{`)); err == nil {
		t.Fatalf("parsePackages() should reject malformed JSON")
	}
}
