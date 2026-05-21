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

func TestCloneStringsReturnsDetachedSlice(t *testing.T) {
	in := []string{"./..."}
	out := cloneStrings(in)
	out[0] = "mutated"
	if in[0] == "mutated" {
		t.Fatalf("cloneStrings leaked input slice")
	}
}

func TestCloneBranchMappingsReturnsDetachedSlice(t *testing.T) {
	mapping, err := NewBranchMapping(BranchMappingSpec{Source: "main", Target: "release"})
	if err != nil {
		t.Fatalf("NewBranchMapping returned error: %v", err)
	}
	in := []BranchMapping{mapping}
	out := CloneBranchMappings(in)
	out[0], _ = NewBranchMapping(BranchMappingSpec{Source: "dev", Target: "dev"})
	if in[0].Source() != "main" {
		t.Fatalf("CloneBranchMappings leaked input slice")
	}
}

func TestCloneModuleNamesReturnsDetachedSlice(t *testing.T) {
	in := []ModuleName{"control"}
	out := CloneModuleNames(in)
	out[0] = "mutated"
	if in[0] == "mutated" {
		t.Fatalf("CloneModuleNames leaked input slice")
	}
}

func TestClonePublishEntriesReturnsDetachedSlice(t *testing.T) {
	entry, err := NewPublishEntry(PublishEntrySpec{Type: "file", From: "go.mod", To: "go.mod"})
	if err != nil {
		t.Fatalf("NewPublishEntry returned error: %v", err)
	}
	in := []PublishEntry{entry}
	out := ClonePublishEntries(in)
	out[0], _ = NewPublishEntry(PublishEntrySpec{Type: "file", From: "README.md", To: "README.md"})
	if in[0].From() != "go.mod" {
		t.Fatalf("ClonePublishEntries leaked input slice")
	}
}
