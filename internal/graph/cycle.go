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
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// Cycle is one dependency cycle in graph module-name form.
type Cycle struct {
	nodes []manifest.ModuleName
}

// Nodes returns the cycle nodes. The returned slice includes the repeated first
// node as the final element when the cycle was discovered from a DFS back edge.
func (c Cycle) Nodes() []manifest.ModuleName { return cloneModuleNames(c.nodes) }

// String returns a stable human-readable cycle path.
func (c Cycle) String() string {
	parts := make([]string, 0, len(c.nodes))
	for _, node := range c.nodes {
		parts = append(parts, node.String())
	}
	return strings.Join(parts, " -> ")
}

// HasCycle reports whether the graph contains at least one dependency cycle.
func (g Graph) HasCycle() bool {
	_, ok := g.FindCycle()
	return ok
}

// Cycles returns discovered dependency cycles in deterministic DFS order.
//
// The current implementation returns at most one cycle. The slice form leaves
// room for future multi-cycle reporting without changing callers.
func (g Graph) Cycles() []Cycle {
	cycle, ok := g.FindCycle()
	if !ok {
		return nil
	}
	return []Cycle{cycle}
}

// FindCycle returns the first dependency cycle found in deterministic DFS order.
func (g Graph) FindCycle() (Cycle, bool) {
	state := make(map[manifest.ModuleName]visitState, len(g.order))
	stack := make([]manifest.ModuleName, 0, len(g.order))
	index := make(map[manifest.ModuleName]int, len(g.order))
	for _, name := range g.order {
		if state[name] != unvisited {
			continue
		}
		if cycle, ok := g.findCycleFrom(name, state, &stack, index); ok {
			return cycle, true
		}
	}
	return Cycle{}, false
}

type visitState uint8

const (
	unvisited visitState = iota
	visiting
	visited
)

// findCycleFrom runs a DFS from name and returns the first back-edge cycle.
func (g Graph) findCycleFrom(
	name manifest.ModuleName,
	state map[manifest.ModuleName]visitState,
	stack *[]manifest.ModuleName,
	index map[manifest.ModuleName]int,
) (Cycle, bool) {
	state[name] = visiting
	index[name] = len(*stack)
	*stack = append(*stack, name)
	for _, dependent := range g.dependents[name] {
		switch state[dependent] {
		case unvisited:
			if cycle, ok := g.findCycleFrom(dependent, state, stack, index); ok {
				return cycle, true
			}
		case visiting:
			start := index[dependent]
			nodes := append([]manifest.ModuleName(nil), (*stack)[start:]...)
			nodes = append(nodes, dependent)
			return Cycle{nodes: nodes}, true
		}
	}
	delete(index, name)
	*stack = (*stack)[:len(*stack)-1]
	state[name] = visited
	return Cycle{}, false
}
