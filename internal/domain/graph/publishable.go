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

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// PublishableSubgraph returns a graph containing only public modules.
//
// Dependencies pointing outside the public subset make the subgraph invalid and
// are returned as validation errors. This keeps accidental publication of a
// public module that depends on an internal or disabled module explicit.
func (g Graph) PublishableSubgraph() (Graph, error) {
	modules, allowed := g.publishableModules()
	if issues := validatePublishableDependencies(modules, allowed); len(issues) > 0 {
		return Graph{}, &ValidationError{Issues: issues}
	}
	return New(modules)
}

// publishableModules returns public modules plus a membership set for validation.
func (g Graph) publishableModules() ([]manifest.Module, map[manifest.ModuleName]struct{}) {
	modules := make([]manifest.Module, 0)
	allowed := map[manifest.ModuleName]struct{}{}
	for _, module := range g.Modules() {
		if module.Publishable() {
			modules = append(modules, module)
			allowed[module.Name()] = struct{}{}
		}
	}
	return modules, allowed
}

// validatePublishableDependencies rejects public modules that depend on hidden modules.
func validatePublishableDependencies(modules []manifest.Module, allowed map[manifest.ModuleName]struct{}) []Issue {
	var issues []Issue
	for _, module := range modules {
		for _, dependency := range module.Dependencies() {
			if _, ok := allowed[dependency.Module()]; !ok {
				issues = append(issues, publishableDependencyIssue(module, dependency.Module()))
			}
		}
	}
	return issues
}

// publishableDependencyIssue creates the domain issue for hidden dependencies.
func publishableDependencyIssue(module manifest.Module, dependency manifest.ModuleName) Issue {
	return Issue{
		Code:    IssueUnknownDependency,
		Module:  module.Name(),
		Message: fmt.Sprintf("publishable module %q depends on non-publishable module %q", module.Name(), dependency),
	}
}
