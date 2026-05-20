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

package process

import "testing"

func TestIsAllowedExitCode_DefaultsToZeroOnly(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		allowed []int
		want    bool
	}{
		{name: "nil accepts zero", code: 0, allowed: nil, want: true},
		{name: "nil rejects non-zero", code: 1, allowed: nil, want: false},
		{name: "empty accepts zero", code: 0, allowed: []int{}, want: true},
		{name: "empty rejects non-zero", code: 1, allowed: []int{}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllowedExitCode(tt.code, tt.allowed); got != tt.want {
				t.Fatalf("IsAllowedExitCode(%d, %v) = %v, want %v", tt.code, tt.allowed, got, tt.want)
			}
		})
	}
}

func TestIsAllowedExitCode_ExplicitAllowedSet(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		allowed []int
		want    bool
	}{
		{name: "zero", code: 0, allowed: []int{0, 1}, want: true},
		{name: "non-zero", code: 1, allowed: []int{0, 1}, want: true},
		{name: "absent", code: 2, allowed: []int{0, 1}, want: false},
		{name: "negative", code: -1, allowed: []int{-1}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllowedExitCode(tt.code, tt.allowed); got != tt.want {
				t.Fatalf("IsAllowedExitCode(%d, %v) = %v, want %v", tt.code, tt.allowed, got, tt.want)
			}
		})
	}
}
