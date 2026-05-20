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

package exec

import (
	"bytes"
	"strings"
)

const redactionMarker = "<redacted>"

// Redactor removes configured sensitive values from rendered command data.
//
// Redaction is literal string replacement. It intentionally avoids regexes so
// secrets are treated as data, not patterns, and so command diagnostics stay
// predictable.
type Redactor struct{ values []string }

// NewRedactor creates a redactor for raw sensitive values. Empty values are ignored.
func NewRedactor(values ...string) Redactor {
	filtered := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			filtered = append(filtered, value)
		}
	}
	return Redactor{values: filtered}
}

// RedactString returns s with every configured sensitive value replaced.
func (r Redactor) RedactString(s string) string {
	for _, value := range r.values {
		s = strings.ReplaceAll(s, value, redactionMarker)
	}
	return s
}

// RedactBytes returns b with every configured sensitive value replaced.
func (r Redactor) RedactBytes(b []byte) []byte {
	out := append([]byte(nil), b...)
	for _, value := range r.values {
		out = bytes.ReplaceAll(out, []byte(value), []byte(redactionMarker))
	}
	return out
}

// RedactSlice returns a detached redacted copy of values.
func (r Redactor) RedactSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = r.RedactString(value)
	}
	return out
}
