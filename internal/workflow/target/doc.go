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

// Package target prepares local target repository worktrees for a publication
// plan.
//
// The package resolves deterministic local worktree paths, verifies that target
// repositories are clean, and records branch publication mappings for later
// workflow stages. It does not copy source files, rewrite go.mod, verify Go
// packages, create commits, tag releases, or push refs.
package target
