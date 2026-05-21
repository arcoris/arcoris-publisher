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

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// validateIndexes checks that lookup maps match the module plan list.
func (v *planValidator) validateIndexes() {
	v.validateModuleIndex()
	v.validatePathIndex()
	v.validateRepositoryIndex()
}

// validateModuleIndex checks the module-name lookup map.
func (v *planValidator) validateModuleIndex() {
	v.validateIndexByModuleName("byModule", v.plan.byModule, v.expected.byModule)
}

// validatePathIndex checks the module-path lookup map.
func (v *planValidator) validatePathIndex() {
	v.validateIndexByModulePath("byPath", v.plan.byPath, v.expected.byPath)
}

// validateRepositoryIndex checks the repository lookup map.
func (v *planValidator) validateRepositoryIndex() {
	v.validateIndexByRepository("byRepo", v.plan.byRepo, v.expected.byRepo)
}

// validateIndexByModuleName compares module-name indexes.
func (v *planValidator) validateIndexByModuleName(index string, actual map[manifest.ModuleName]int, expected map[manifest.ModuleName]int) {
	if len(actual) != len(expected) {
		v.addInvalidIndexIssue(index, fmt.Sprintf("index %s contains %d entries, want %d", index, len(actual), len(expected)))
	}
	for key, want := range expected {
		got, ok := actual[key]
		if !ok {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s is missing module %q", index, key))
			continue
		}
		if got != want {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s maps module %q to %d, want %d", index, key, got, want))
		}
	}
	for key := range actual {
		if _, ok := expected[key]; !ok {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s contains unexpected module %q", index, key))
		}
	}
}

// validateIndexByModulePath compares module-path indexes.
func (v *planValidator) validateIndexByModulePath(index string, actual map[manifest.ModulePath]int, expected map[manifest.ModulePath]int) {
	if len(actual) != len(expected) {
		v.addInvalidIndexIssue(index, fmt.Sprintf("index %s contains %d entries, want %d", index, len(actual), len(expected)))
	}
	for key, want := range expected {
		got, ok := actual[key]
		if !ok {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s is missing module path %q", index, key))
			continue
		}
		if got != want {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s maps module path %q to %d, want %d", index, key, got, want))
		}
	}
	for key := range actual {
		if _, ok := expected[key]; !ok {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s contains unexpected module path %q", index, key))
		}
	}
}

// validateIndexByRepository compares repository indexes.
func (v *planValidator) validateIndexByRepository(index string, actual map[manifest.RepositoryRef]int, expected map[manifest.RepositoryRef]int) {
	if len(actual) != len(expected) {
		v.addInvalidIndexIssue(index, fmt.Sprintf("index %s contains %d entries, want %d", index, len(actual), len(expected)))
	}
	for key, want := range expected {
		got, ok := actual[key]
		if !ok {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s is missing repository %q", index, key))
			continue
		}
		if got != want {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s maps repository %q to %d, want %d", index, key, got, want))
		}
	}
	for key := range actual {
		if _, ok := expected[key]; !ok {
			v.addInvalidIndexIssue(index, fmt.Sprintf("index %s contains unexpected repository %q", index, key))
		}
	}
}
