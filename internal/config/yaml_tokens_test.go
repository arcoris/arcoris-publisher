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

func TestSplitYAMLKeyValue(t *testing.T) {
	key, raw, hasValue, err := splitYAMLKeyValue(yamlLine{
		number: 1,
		text:   "name: arcoris",
	})
	if err != nil {
		t.Fatal(err)
	}
	if key != "name" || raw != "arcoris" || !hasValue {
		t.Fatalf("unexpected split: %q %q %v", key, raw, hasValue)
	}
}

func TestSplitYAMLKeyValueRejectsEmptyKey(t *testing.T) {
	_, _, _, err := splitYAMLKeyValue(yamlLine{
		number: 1,
		text:   ": value",
	})
	if err == nil {
		t.Fatal("expected empty key error")
	}
}

func TestSplitYAMLKeyValueRejectsMissingColon(t *testing.T) {
	_, _, _, err := splitYAMLKeyValue(yamlLine{
		number: 1,
		text:   "name",
	})
	if err == nil {
		t.Fatal("expected key-value error")
	}
}

func TestValidateListItemLine(t *testing.T) {
	ok, err := validateListItemLine(yamlLine{indent: 2, text: "- name"}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected list item")
	}
}

func TestValidateListItemLineRejectsUnexpectedIndent(t *testing.T) {
	_, err := validateListItemLine(yamlLine{number: 1, indent: 4, text: "- name"}, 2)
	if err == nil {
		t.Fatal("expected indentation error")
	}
}
