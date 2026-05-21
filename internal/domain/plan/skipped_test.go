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

package plan

import "testing"

func TestSkippedModules(t *testing.T) {
	skipped := testPlan(t).SkippedModules()
	if len(skipped) != 2 {
		t.Fatalf("len(SkippedModules()) = %d, want 2", len(skipped))
	}
	if skipped[0].Name() != moduleName(t, "internal-tools") || skipped[0].Reason() != SkipInternal {
		t.Fatalf("unexpected first skipped module: %#v", skipped[0])
	}
	if skipped[1].Name() != moduleName(t, "old-module") || skipped[1].Reason() != SkipDisabled {
		t.Fatalf("unexpected second skipped module: %#v", skipped[1])
	}
	if skipped[0].Module().Name() != skipped[0].Name() {
		t.Fatalf("Module().Name() = %q", skipped[0].Module().Name())
	}
}
