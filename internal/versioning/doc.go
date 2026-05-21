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

// Package versioning assigns publication versions to resolved ARCORIS Publisher
// modules.
//
// The package consumes the resolved publication model, registry indexes, and the
// dependency graph to produce deterministic module version assignments and
// direct internal dependency requirements. It deliberately does not read
// manifests, inspect Git tags, rewrite go.mod files, access the filesystem, call
// Git, or execute publication workflows.
package versioning
