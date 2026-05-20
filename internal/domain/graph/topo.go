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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// TopologicalOrder returns modules in dependency-first order.
//
// The returned order is deterministic: when several modules are simultaneously
// ready, declaration order from the manifest input is used as the tie-breaker.
func (g Graph) TopologicalOrder() ([]manifest.ModuleName, error) {
	if err := g.validateAcyclicForOrdering(); err != nil {
		return nil, err
	}
	sorter := newTopoSorter(g)
	order := sorter.sort()
	if len(order) != len(g.order) {
		return nil, cycleValidationError([]Cycle{})
	}
	return order, nil
}

// PublishOrder is an alias for TopologicalOrder.
func (g Graph) PublishOrder() ([]manifest.ModuleName, error) {
	return g.TopologicalOrder()
}

// validateAcyclicForOrdering converts discovered cycles into ordering errors.
func (g Graph) validateAcyclicForOrdering() error {
	if cycles := g.Cycles(); len(cycles) > 0 {
		return cycleValidationError(cycles)
	}
	return nil
}

// cycleValidationError builds a graph validation error from cycle data.
func cycleValidationError(cycles []Cycle) error {
	if len(cycles) == 0 {
		return &ValidationError{Issues: []Issue{{Code: IssueCycle, Message: "dependency cycle detected"}}}
	}
	issues := make([]Issue, 0, len(cycles))
	for _, cycle := range cycles {
		issues = append(issues, cycleIssue(cycle))
	}
	return &ValidationError{Issues: issues}
}

// topoSorter owns the mutable Kahn algorithm state for one ordering pass.
type topoSorter struct {
	graph    Graph
	inDegree map[manifest.ModuleName]int
	ready    []manifest.ModuleName
	result   []manifest.ModuleName
}

// newTopoSorter prepares in-degree and ready-queue state.
func newTopoSorter(graph Graph) topoSorter {
	sorter := topoSorter{
		graph:    graph,
		inDegree: make(map[manifest.ModuleName]int, len(graph.order)),
		ready:    make([]manifest.ModuleName, 0, len(graph.order)),
		result:   make([]manifest.ModuleName, 0, len(graph.order)),
	}
	sorter.initialize()
	return sorter
}

// initialize computes initial in-degree and declaration-ordered ready nodes.
func (s *topoSorter) initialize() {
	for _, name := range s.graph.order {
		s.inDegree[name] = len(s.graph.dependencies[name])
	}
	for _, name := range s.graph.order {
		if s.inDegree[name] == 0 {
			s.ready = append(s.ready, name)
		}
	}
}

// sort drains the ready queue and returns the dependency-first order.
func (s *topoSorter) sort() []manifest.ModuleName {
	for len(s.ready) > 0 {
		current := s.popReady()
		s.result = append(s.result, current)
		s.releaseDependents(current)
	}
	return s.result
}

// popReady removes the first declaration-ordered ready module.
func (s *topoSorter) popReady() manifest.ModuleName {
	current := s.ready[0]
	s.ready = s.ready[1:]
	return current
}

// releaseDependents reduces in-degree for modules unlocked by current.
func (s *topoSorter) releaseDependents(current manifest.ModuleName) {
	for _, dependent := range s.graph.dependents[current] {
		s.inDegree[dependent]--
		if s.inDegree[dependent] == 0 {
			s.ready = insertByDeclarationOrder(s.ready, dependent, s.graph.order)
		}
	}
}

// insertByDeclarationOrder inserts value into an already declaration-ordered slice.
func insertByDeclarationOrder(values []manifest.ModuleName, value manifest.ModuleName, order []manifest.ModuleName) []manifest.ModuleName {
	valueIndex := declarationIndex(value, order)
	insertAt := len(values)
	for i, existing := range values {
		if declarationIndex(existing, order) > valueIndex {
			insertAt = i
			break
		}
	}
	values = append(values, "")
	copy(values[insertAt+1:], values[insertAt:])
	values[insertAt] = value
	return values
}

// declarationIndex returns the manifest declaration index or len(order) if missing.
func declarationIndex(value manifest.ModuleName, order []manifest.ModuleName) int {
	for i, candidate := range order {
		if candidate == value {
			return i
		}
	}
	return len(order)
}
