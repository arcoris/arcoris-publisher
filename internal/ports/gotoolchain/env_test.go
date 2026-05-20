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

func TestEnvResultValue(t *testing.T) {
	result := EnvResult{Values: map[string]string{"GOWORK": "off", "GOFLAGS": ""}}

	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "present", key: "GOWORK", want: "off"},
		{name: "present empty", key: "GOFLAGS", want: ""},
		{name: "missing", key: "GOPATH", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := result.Value(tt.key); got != tt.want {
				t.Fatalf("Value(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}

	if got := (EnvResult{}).Value("GOWORK"); got != "" {
		t.Fatalf("nil Values Value() = %q, want empty", got)
	}
}

func TestEnvResultHasValue(t *testing.T) {
	result := EnvResult{Values: map[string]string{"GOFLAGS": ""}}

	if !result.HasValue("GOFLAGS") {
		t.Fatalf("HasValue should distinguish present empty values from missing values")
	}
	if result.HasValue("GOWORK") {
		t.Fatalf("HasValue should reject missing keys")
	}
	if (EnvResult{}).HasValue("GOWORK") {
		t.Fatalf("HasValue should reject nil value maps")
	}
}
