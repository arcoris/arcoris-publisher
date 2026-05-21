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

package manifest_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestParseRepositoryRefAcceptsOwnerName(t *testing.T) {
	got, err := manifest.ParseRepositoryRef("arcoris/publisher")
	if err != nil {
		t.Fatalf("ParseRepositoryRef returned error: %v", err)
	}
	if got.String() != "arcoris/publisher" {
		t.Fatalf("unexpected repository ref: %q", got.String())
	}
}

func TestParseRepositoryRefRejectsUnsafeForms(t *testing.T) {
	for _, value := range []string{
		"",
		"arcoris",
		"arcoris/",
		"/publisher",
		"arc/or/is",
		"arc oris/publisher",
		".arc/publisher",
		"arc/../publisher",
	} {
		if _, err := manifest.ParseRepositoryRef(value); err == nil {
			t.Fatalf("ParseRepositoryRef(%q) returned nil error", value)
		}
	}
}
