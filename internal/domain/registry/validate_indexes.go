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

package registry

import (
	"fmt"
	"slices"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// indexValidator compares a registry with a freshly rebuilt expected registry.
type indexValidator struct {
	actual   Registry
	expected Registry
	issues   []Issue
}

// newIndexValidator prepares a deterministic structural registry validator.
func newIndexValidator(actual Registry, expected Registry) indexValidator {
	return indexValidator{actual: actual, expected: expected}
}

// validate returns nil when every internal index matches the module slice.
func (v indexValidator) validate() error {
	validateScalarIndex(&v, "byName", v.actual.byName, v.expected.byName)
	validateScalarIndex(&v, "byModulePath", v.actual.byModulePath, v.expected.byModulePath)
	validateScalarIndex(&v, "bySourceDir", v.actual.bySourceDir, v.expected.bySourceDir)
	validateScalarIndex(&v, "byRepository", v.actual.byRepository, v.expected.byRepository)
	v.validateVisibilityIndex()
	v.validateBranchIndex()

	if len(v.issues) == 0 {
		return nil
	}
	return &ValidationError{Issues: v.issues}
}

// validateScalarIndex checks one unique-key index whose values are module positions.
func validateScalarIndex[K comparable](v *indexValidator, name string, actual map[K]int, expected map[K]int) {
	if len(actual) != len(expected) {
		v.addInvalidIndexIssue(name, "index %s contains %d entries, want %d", name, len(actual), len(expected))
	}

	for key, want := range expected {
		got, ok := actual[key]
		switch {
		case !ok:
			v.addInvalidIndexIssue(name, "index %s is missing key %v", name, key)
		case got != want:
			v.addInvalidIndexIssue(name, "index %s maps key %v to module index %d, want %d", name, key, got, want)
		}
	}

	for key := range actual {
		if _, ok := expected[key]; !ok {
			v.addInvalidIndexIssue(name, "index %s contains unexpected key %v", name, key)
		}
	}
}

// validateVisibilityIndex checks visibility buckets and their declaration order.
func (v *indexValidator) validateVisibilityIndex() {
	actual := v.actual.byVisibility
	expected := v.expected.byVisibility
	if len(actual) != len(expected) {
		v.addInvalidIndexIssue("byVisibility", "index byVisibility contains %d buckets, want %d", len(actual), len(expected))
	}

	for visibility, want := range expected {
		got, ok := actual[visibility]
		switch {
		case !ok:
			v.addInvalidIndexIssue("byVisibility", "index byVisibility is missing visibility %q", visibility)
		case !slices.Equal(got, want):
			v.addInvalidIndexIssue("byVisibility", "index byVisibility maps visibility %q to %v, want %v", visibility, got, want)
		}
	}

	for visibility := range actual {
		if _, ok := expected[visibility]; !ok {
			v.addInvalidIndexIssue("byVisibility", "index byVisibility contains unexpected visibility %q", visibility)
		}
	}
}

// validateBranchIndex checks each module's source-branch mapping index.
func (v *indexValidator) validateBranchIndex() {
	actual := v.actual.byBranch
	expected := v.expected.byBranch
	if len(actual) != len(expected) {
		v.addInvalidIndexIssue("byBranch", "index byBranch contains %d modules, want %d", len(actual), len(expected))
	}

	for module, want := range expected {
		got, ok := actual[module]
		switch {
		case !ok:
			v.addInvalidIndexIssue("byBranch", "index byBranch is missing module %q", module)
		default:
			v.validateBranchMappings(module, got, want)
		}
	}

	for module := range actual {
		if _, ok := expected[module]; !ok {
			v.addInvalidIndexIssue("byBranch", "index byBranch contains unexpected module %q", module)
		}
	}
}

// validateBranchMappings checks one module's source-branch map.
func (v *indexValidator) validateBranchMappings(
	module manifest.ModuleName,
	actual map[manifest.BranchName]manifest.BranchMapping,
	expected map[manifest.BranchName]manifest.BranchMapping,
) {
	indexName := fmt.Sprintf("byBranch[%s]", module)
	if len(actual) != len(expected) {
		v.addInvalidIndexIssue(indexName, "index %s contains %d branches, want %d", indexName, len(actual), len(expected))
	}

	for source, want := range expected {
		got, ok := actual[source]
		switch {
		case !ok:
			v.addInvalidIndexIssue(indexName, "index %s is missing source branch %q", indexName, source)
		case got != want:
			v.addInvalidIndexIssue(indexName, "index %s maps source branch %q to %q, want %q", indexName, source, got.Target(), want.Target())
		}
	}

	for source := range actual {
		if _, ok := expected[source]; !ok {
			v.addInvalidIndexIssue(indexName, "index %s contains unexpected source branch %q", indexName, source)
		}
	}
}

// addInvalidIndexIssue records a structural registry inconsistency.
func (v *indexValidator) addInvalidIndexIssue(index string, format string, args ...any) {
	v.issues = append(v.issues, invalidIndexIssue(index, fmt.Sprintf(format, args...)))
}
