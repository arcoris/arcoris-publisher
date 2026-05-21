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

func TestValidateSimpleNameAcceptsLowerKebabIdentifiers(t *testing.T) {
	if err := validateSimpleName("field", "alpha-1"); err != nil {
		t.Fatalf("validateSimpleName returned error: %v", err)
	}
}

func TestValidateSimpleNameRejectsMissingWhitespaceAndPatternMismatches(t *testing.T) {
	for _, value := range []string{"", " alpha", "Alpha", "alpha_1"} {
		if err := validateSimpleName("field", value); err == nil {
			t.Fatalf("validateSimpleName(%q) returned nil error", value)
		}
	}
}
