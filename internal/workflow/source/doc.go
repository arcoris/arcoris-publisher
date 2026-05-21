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

// Package source inspects the source checkout for an executable publication
// plan.
//
// The package is the first runtime workflow stage after planning. It validates
// that the current source repository can satisfy the plan, captures Git source
// state, checks explicit publish entries, and returns an immutable source
// snapshot used by later target, construction, module-file, verification, and
// publication stages.
//
// This package does not load manifests, build registries, build dependency
// graphs, assign versions, create target worktrees, copy files to targets,
// rewrite go.mod files, run Go verification, commit, tag, or push refs.
package source
