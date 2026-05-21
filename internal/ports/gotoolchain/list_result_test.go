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

func TestListResultHasPackages(t *testing.T) {
	tests := []struct {
		name   string
		result ListResult
		want   bool
	}{
		{name: "zero", result: ListResult{}, want: false},
		{name: "empty slice", result: ListResult{Packages: []Package{}}, want: false},
		{name: "package", result: ListResult{Packages: []Package{{ImportPath: "example.com/pkg"}}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.HasPackages(); got != tt.want {
				t.Fatalf("HasPackages() = %v, want %v", got, tt.want)
			}
		})
	}
}
