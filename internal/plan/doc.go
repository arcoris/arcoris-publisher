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

// Package plan builds executable publication plans from resolved ARCORIS
// Publisher manifests, registry indexes, dependency graphs, and version
// assignments.
//
// A plan is an immutable description of what should be published. It contains
// module order, source and target routing, explicit publish entries, assigned
// versions, dependency requirements, branch mappings, and effective
// verification policy.
//
// This package deliberately does not read manifest files, access the
// filesystem, call Git, rewrite go.mod files, run the Go toolchain, create
// commits, tag repositories, or push refs. Runtime workflow packages consume a
// Plan and perform those side effects later.
package plan
