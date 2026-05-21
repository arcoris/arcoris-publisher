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

// Package graph builds dependency topology for resolved ARCORIS Publisher modules.
//
// The package consumes the registry built from the resolved publication model and
// provides deterministic dependency traversal, cycle discovery, topological
// ordering, publication ordering, and impact analysis.
//
// Graph edges are directed from dependency to dependent. An edge A -> B means
// that module B depends on module A and A must be published before B.
//
// This package does not load manifests, decode YAML, assign versions, build
// publication plans, access the filesystem, call Git, invoke the Go toolchain,
// or execute workflows.
package graph
