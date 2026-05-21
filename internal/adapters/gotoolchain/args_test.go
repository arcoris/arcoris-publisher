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

func TestDefaultPatterns(t *testing.T) {
	assertStringSlice(t, defaultPatterns(nil), []string{"./..."})

	patterns := []string{"./cmd"}
	got := defaultPatterns(patterns)
	patterns[0] = "./mutated"
	assertStringSlice(t, got, []string{"./cmd"})
}

func TestAppendBuildTags(t *testing.T) {
	assertStringSlice(t, appendBuildTags([]string{"test"}, nil), []string{"test"})
	assertStringSlice(t, appendBuildTags([]string{"test"}, []string{"integration", "linux"}), []string{"test", "-tags", "integration,linux"})
}
