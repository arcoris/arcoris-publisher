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

// Package provenance builds deterministic publication metadata.
//
// The package is deliberately transform-only. It does not read manifests,
// inspect Git, touch the filesystem, build plans, copy files, or publish refs.
// Callers provide resolved runtime objects, and provenance renders stable JSON
// payloads, commit trailers, and explicit-projection hashes without leaking
// local absolute source or target paths.
package provenance
