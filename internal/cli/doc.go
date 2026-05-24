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

// Package cli implements command routing, flag parsing, report rendering, and
// process exit-code decisions for the ARCORIS Publisher command-line interface.
//
// The package is deliberately thin. It delegates high-level use cases to
// internal/app and delegates output formatting to internal/report. It does not
// load manifests directly, build dependency graphs, execute workflow stages,
// inspect Git repositories, copy files, rewrite go.mod files, or publish refs.
// Concrete infrastructure adapters are wired outside the workflow packages and
// passed in through application dependencies.
package cli
