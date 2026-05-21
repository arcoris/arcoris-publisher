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

// Package modulefile rewrites constructed Go module files so target repositories
// are publishable as standalone modules.
//
// The package consumes plan dependency requirements and target worktrees. It
// updates module directives, direct internal requirements, and local replace
// directives. It does not compute dependency versions, run go mod tidy, verify
// packages, commit, tag, or push repositories.
package modulefile
