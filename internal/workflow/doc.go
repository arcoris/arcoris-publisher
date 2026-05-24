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

// Package workflow orchestrates publication workflow stages over an already
// built plan.
//
// The runner deliberately does not load manifests, build registries, compute
// dependency graphs, assign versions, or construct plans. Those are application
// responsibilities. This package wires stage services in order and preserves
// each stage boundary: source inspection, target preparation, construction,
// module-file rewrite, verification, and optional publication.
package workflow
