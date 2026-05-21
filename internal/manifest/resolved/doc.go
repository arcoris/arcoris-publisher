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

// Package resolved binds staging and module manifests into an effective
// publication model consumed by registry, graph, plan, and workflow packages.
//
// The package applies built-in defaults, top-level defaults, staging module
// overrides, and module-level overrides. Downstream packages should consume
// PublicationSet instead of raw manifest resources.
package resolved
