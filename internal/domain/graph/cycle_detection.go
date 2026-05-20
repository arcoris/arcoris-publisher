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

// Acyclic reports whether the graph contains no dependency cycles.
func (g Graph) Acyclic() bool {
	return len(g.Cycles()) == 0
}

// Cycles returns dependency cycles in deterministic order.
func (g Graph) Cycles() []Cycle {
	finder := newCycleFinder(g)
	finder.find()
	return finder.cycles
}

// cycleFinder owns DFS state for one cycle-detection pass.
type cycleFinder struct {
	graph      Graph
	state      map[manifest.ModuleName]visitState
	stack      []manifest.ModuleName
	stackIndex map[manifest.ModuleName]int
	seen       map[string]struct{}
	cycles     []Cycle
}

// newCycleFinder prepares DFS state sized for the graph.
func newCycleFinder(graph Graph) cycleFinder {
	return cycleFinder{
		graph:      graph,
		state:      make(map[manifest.ModuleName]visitState, len(graph.order)),
		stack:      make([]manifest.ModuleName, 0, len(graph.order)),
		stackIndex: make(map[manifest.ModuleName]int, len(graph.order)),
		seen:       map[string]struct{}{},
		cycles:     make([]Cycle, 0),
	}
}

// find visits every disconnected component in declaration order.
func (f *cycleFinder) find() {
	for _, module := range f.graph.order {
		if f.state[module] == visitUnseen {
			f.visit(module)
		}
	}
}

// visit performs depth-first traversal over dependency edges.
func (f *cycleFinder) visit(module manifest.ModuleName) {
	f.enter(module)
	for _, dependency := range f.graph.dependencies[module] {
		switch f.state[dependency] {
		case visitUnseen:
			f.visit(dependency)
		case visitActive:
			f.recordCycle(dependency)
		}
	}
	f.leave(module)
}

// enter marks module active and records its stack position.
func (f *cycleFinder) enter(module manifest.ModuleName) {
	f.state[module] = visitActive
	f.stackIndex[module] = len(f.stack)
	f.stack = append(f.stack, module)
}

// leave marks module done and removes it from the active stack.
func (f *cycleFinder) leave(module manifest.ModuleName) {
	f.stack = f.stack[:len(f.stack)-1]
	delete(f.stackIndex, module)
	f.state[module] = visitDone
}

// recordCycle extracts, normalizes, and deduplicates an active-stack cycle.
func (f *cycleFinder) recordCycle(start manifest.ModuleName) {
	cycle := f.stackCycle(start)
	key := cycle.String()
	if _, exists := f.seen[key]; exists {
		return
	}
	f.seen[key] = struct{}{}
	f.cycles = append(f.cycles, cycle)
}

// stackCycle returns the current stack slice from start back to start.
func (f *cycleFinder) stackCycle(start manifest.ModuleName) Cycle {
	startIndex := f.stackIndex[start]
	nodes := append([]manifest.ModuleName(nil), f.stack[startIndex:]...)
	nodes = append(nodes, start)
	return normalizeCycle(nodes)
}

// visitState stores DFS color for cycle detection.
type visitState uint8

const (
	visitUnseen visitState = iota
	visitActive
	visitDone
)
