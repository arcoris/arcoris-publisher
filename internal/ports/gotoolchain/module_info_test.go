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

func TestModuleInfoHasReplace(t *testing.T) {
	tests := []struct {
		name   string
		module ModuleInfo
		want   bool
	}{
		{name: "zero", module: ModuleInfo{}, want: false},
		{name: "replace", module: ModuleInfo{Replace: &ModuleInfo{Path: "example.com/old"}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.module.HasReplace(); got != tt.want {
				t.Fatalf("HasReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}
