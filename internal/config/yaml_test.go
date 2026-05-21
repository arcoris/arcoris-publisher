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
	"encoding/json"
	"strings"
	"testing"
)

func TestYAMLToJSONHandlesCommentsQuotesBooleansAndNumbers(t *testing.T) {
	data := []byte(`name: "arcoris" # comment
single: 'it''s ok'
enabled: false
count: 7
ratio: 1.5
emptyList: []
emptyMap: {}
nullable: null
hash: "value # not comment"
`)
	jsonData, err := yamlToJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["name"] != "arcoris" || decoded["single"] != "it's ok" || decoded["enabled"] != false {
		t.Fatalf("unexpected decoded map: %#v", decoded)
	}
	if decoded["hash"] != "value # not comment" {
		t.Fatalf("comment inside quoted string was stripped: %#v", decoded["hash"])
	}
}

func TestYAMLToJSONHandlesNestedListItems(t *testing.T) {
	data := []byte(`modules:
  - name: foundation
    sourceDir: src/foundation
    branches:
      - source: main
        target: main
  - name: control
    sourceDir: src/control
`)
	jsonData, err := yamlToJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string][]map[string]any
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded["modules"]) != 2 || decoded["modules"][0]["name"] != "foundation" {
		t.Fatalf("unexpected decoded modules: %#v", decoded)
	}
}

func TestYAMLToJSONHandlesNestedDashItems(t *testing.T) {
	data := []byte(`items:
  -
    name: foundation
`)
	jsonData, err := yamlToJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	var decoded map[string][]map[string]any
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["items"][0]["name"] != "foundation" {
		t.Fatalf("unexpected decoded items: %#v", decoded)
	}
}

func TestYAMLToJSONHandlesTopLevelList(t *testing.T) {
	jsonData, err := yamlToJSON([]byte("- one\n- two\n"))
	if err != nil {
		t.Fatal(err)
	}

	var decoded []string
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 2 || decoded[0] != "one" {
		t.Fatalf("unexpected decoded list: %#v", decoded)
	}
}

func TestYAMLToJSONHandlesEmptyDocument(t *testing.T) {
	jsonData, err := yamlToJSON([]byte("\n# comment\n"))
	if err != nil {
		t.Fatal(err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 0 {
		t.Fatalf("unexpected decoded document: %#v", decoded)
	}
}

func TestYAMLToJSONRejectsTabs(t *testing.T) {
	if _, err := yamlToJSON([]byte("key:\n\tchild: value\n")); err == nil {
		t.Fatal("expected tab rejection")
	}
}

func TestYAMLToJSONRejectsUnexpectedIndentation(t *testing.T) {
	if _, err := yamlToJSON([]byte("key: value\n  child: value\n")); err == nil {
		t.Fatal("expected indentation error")
	}
}

func TestYAMLToJSONRejectsMalformedKey(t *testing.T) {
	if _, err := yamlToJSON([]byte("bad key: value\n")); err == nil {
		t.Fatal("expected malformed key error")
	}
}

func TestYAMLToJSONRejectsUnterminatedQuote(t *testing.T) {
	_, err := yamlToJSON([]byte("key: \"unterminated\n"))
	if err == nil {
		t.Fatal("expected quote error")
	}

	hasUnterminated := strings.Contains(err.Error(), "unterminated")
	hasQuoted := strings.Contains(err.Error(), "quoted")
	if !hasUnterminated && !hasQuoted {
		t.Fatalf("expected quote error, got %v", err)
	}
}
