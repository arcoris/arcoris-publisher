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

// Package manifest contains the shared value types, policies, diagnostics, and
// validation primitives used by ARCORIS Publisher manifests.
//
// The root package intentionally contains only common concepts. The top-level
// arcpub.yaml resource lives in package staging. The module-level
// arcpub.module.yaml resource lives in package module. The package resolved
// binds those two resource kinds into an effective publication model.
//
// This package does not read files, decode YAML, call Git, access the
// filesystem, or execute publication workflows.
package manifest
