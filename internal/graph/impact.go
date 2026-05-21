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

package graph

import "arcoris.dev/arcoris-publisher/internal/manifest"

// AffectedBy returns changed modules and all transitive dependents.
//
// The result is returned in topological order so dependencies appear before the
// modules that depend on them. Unknown input names are reported as validation
// issues because affected-set operations are normally user-facing.
func (g Graph) AffectedBy(changed ...manifest.ModuleName) ([]manifest.ModuleName, error) {
	selected, err := g.validateInputNames("changed", changed)
	if err != nil {
		return nil, err
	}
	seen := make(map[manifest.ModuleName]struct{}, len(selected))
	for _, name := range selected {
		seen[name] = struct{}{}
		g.collectDependents(name, seen)
	}
	return g.filterByTopologicalOrder(seen), nil
}

// PublishClosure returns selected public modules and their transitive dependencies.
//
// This is the dependency closure needed to publish selected modules safely.
func (g Graph) PublishClosure(selected ...manifest.ModuleName) ([]manifest.ModuleName, error) {
	input, err := g.validateInputNames("selected", selected)
	if err != nil {
		return nil, err
	}
	seen := make(map[manifest.ModuleName]struct{}, len(input))
	for _, name := range input {
		seen[name] = struct{}{}
		g.collectDependencies(name, seen)
	}
	return g.filterByTopologicalOrder(seen), nil
}

// validateInputNames rejects unknown graph nodes and deduplicates valid inputs
// in caller-provided order.
func (g Graph) validateInputNames(
	path string,
	names []manifest.ModuleName,
) ([]manifest.ModuleName, error) {
	collector := newIssueCollector()
	seen := make(map[manifest.ModuleName]struct{}, len(names))
	out := make([]manifest.ModuleName, 0, len(names))
	for _, name := range names {
		if !g.Contains(name) {
			collector.Add(
				IssueUnknownNode,
				inputIssuePath(path, len(names)),
				"unknown graph node %q",
				name,
			)
			continue
		}

		if _, ok := seen[name]; ok {
			continue
		}

		seen[name] = struct{}{}
		out = append(out, name)
	}

	if err := collector.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// inputIssuePath marks repeated user input as a collection-level location.
func inputIssuePath(path string, inputCount int) string {
	if inputCount <= 1 {
		return path
	}

	return path + "[]"
}
