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

// Package runtime assembles production dependencies for ARCORIS Publisher.
//
// The package is the composition root for command-line binaries. It wires local
// filesystem, process, Git CLI, Go toolchain, manifest loading, application use
// cases, and command routing together without moving infrastructure concerns
// into workflow packages.
//
// Runtime deliberately does not parse command-line flags, execute workflow
// stages directly, render reports, or implement publication policy. Command
// routing belongs to internal/cli, high-level use cases belong to internal/app,
// and workflow stages remain port-oriented.
package runtime
