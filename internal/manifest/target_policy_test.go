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

func TestTargetPolicyRemoteTemplate(t *testing.T) {
	raw := "https://github.com/{repository}.git"
	policy, err := NewTargetPolicy(TargetSpec{RemoteTemplate: &raw})
	if err != nil {
		t.Fatal(err)
	}
	template, ok := policy.RemoteTemplate()
	if !ok || template.String() != raw {
		t.Fatalf("RemoteTemplate() = %q, %v", template, ok)
	}
}

func TestTargetPolicyAllowsOmittedTemplate(t *testing.T) {
	policy, err := NewTargetPolicy(TargetSpec{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := policy.RemoteTemplate(); ok {
		t.Fatal("RemoteTemplate() ok = true")
	}
}
