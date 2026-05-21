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

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// TopologicalOrder returns all graph nodes in dependency-before-dependent order.
//
// The method uses Kahn's algorithm with registry declaration order as the queue
// tie-breaker. If the graph contains a dependency cycle, a ValidationError with
// IssueDependencyCycle is returned.
func (g Graph) TopologicalOrder() ([]manifest.ModuleName, error) {
	order, err := g.topologicalOrder(nil)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// PublishOrder returns public modules in dependency-before-dependent order.
//
// Internal modules may participate in topology but are omitted from the returned
// order. Disabled modules are not present in the graph.
func (g Graph) PublishOrder() ([]manifest.ModuleName, error) {
	filter := func(node Node) bool { return node.Publishable() }
	order, err := g.topologicalOrder(filter)
	if err != nil {
		return nil, err
	}
	return order, nil
}

type nodeFilter func(Node) bool

// topologicalOrder returns dependency-before-dependent order and optionally
// filters emitted nodes while still traversing the complete graph.
func (g Graph) topologicalOrder(filter nodeFilter) ([]manifest.ModuleName, error) {
	inDegree := make(map[manifest.ModuleName]int, len(g.order))
	for _, name := range g.order {
		inDegree[name] = len(g.dependencies[name])
	}
	queue := make([]manifest.ModuleName, 0, len(g.order))
	for _, name := range g.order {
		if inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}
	visited := 0
	ordered := make([]manifest.ModuleName, 0, len(g.order))
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		visited++
		if filter == nil || filter(g.nodes[name]) {
			ordered = append(ordered, name)
		}
		for _, dependent := range g.dependents[name] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}
	if visited != len(g.order) {
		cycle, ok := g.FindCycle()
		message := "dependency cycle detected"
		if ok {
			message = fmt.Sprintf("dependency cycle detected: %s", cycle.String())
		}
		return nil, dependencyCycleError(message)
	}
	return ordered, nil
}

// dependencyCycleError builds the common cycle validation error used by ordering
// and explicit acyclicity checks.
func dependencyCycleError(message string) error {
	return &ValidationError{
		Issues: []Issue{
			{
				Code:    IssueDependencyCycle,
				Path:    "dependencies",
				Message: message,
			},
		},
	}
}
