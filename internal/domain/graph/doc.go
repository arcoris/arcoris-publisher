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

// Package graph defines the dependency graph used by ARCORIS Publisher domain
// planning.
//
// The package works with already validated manifest modules. It does not load
// configuration files, access the filesystem, call Git, invoke the Go toolchain,
// or execute publishing workflows. Edges are interpreted in publication order:
// when module A depends on module B, the graph stores the ordering edge B -> A.
// This makes topological order directly usable as a safe publication order.
package graph
