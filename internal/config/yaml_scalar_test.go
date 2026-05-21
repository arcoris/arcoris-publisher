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

func TestParseYAMLScalarParsesEmptyCollectionsAndNull(t *testing.T) {
	cases := []string{"[]", "{}", "null", "~"}
	for _, raw := range cases {
		if _, err := parseYAMLScalar(raw); err != nil {
			t.Fatalf("parseYAMLScalar(%q): %v", raw, err)
		}
	}
}

func TestParseYAMLScalarRejectsUnsupportedInlineCollection(t *testing.T) {
	if _, err := parseYAMLScalar("[value]"); err == nil {
		t.Fatal("expected inline collection error")
	}
}

func TestParseSingleQuotedYAMLScalarUnescapesQuotes(t *testing.T) {
	got, err := parseSingleQuotedYAMLScalar("'it''s ok'")
	if err != nil {
		t.Fatal(err)
	}
	if got != "it's ok" {
		t.Fatalf("got %q", got)
	}
}

func TestParseSingleQuotedYAMLScalarRejectsUnterminated(t *testing.T) {
	if _, err := parseSingleQuotedYAMLScalar("'unterminated"); err == nil {
		t.Fatal("expected unterminated quote error")
	}
}

func TestParseYAMLNumberRejectsNonNumber(t *testing.T) {
	if _, ok := parseYAMLNumber("v1.2.3"); ok {
		t.Fatal("expected non-number")
	}
}
