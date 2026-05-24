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

package report

import (
	"errors"
	"testing"
)

func TestParseFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  Format
		ok    bool
	}{
		{name: "default", input: "", want: FormatText, ok: true},
		{name: "text", input: "text", want: FormatText, ok: true},
		{name: "json", input: "json", want: FormatJSON, ok: true},
		{name: "trim and lower", input: " JSON ", want: FormatJSON, ok: true},
		{name: "invalid", input: "yaml", ok: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseFormat(tt.input)
			if tt.ok && err != nil {
				t.Fatalf("ParseFormat() error = %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("ParseFormat() error = nil")
			}
			if got != tt.want {
				t.Fatalf("ParseFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseFormatUnsupportedReturnsTypedError(t *testing.T) {
	t.Parallel()

	_, err := ParseFormat("yaml")
	if err == nil {
		t.Fatal("ParseFormat() error = nil")
	}

	var reportErr *Error
	if !errors.As(err, &reportErr) || reportErr.Code != CodeUnsupportedFormat {
		t.Fatalf("ParseFormat() error = %v", err)
	}
}
