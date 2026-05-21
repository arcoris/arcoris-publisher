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

// Package versioning defines domain objects for assigning publish versions to
// ARCORIS Publisher modules.
//
// The package is intentionally policy-oriented rather than Git-oriented. It does
// not inspect tags, read commits, call Git, access the filesystem, or execute the
// Go toolchain. Callers provide already known release or snapshot inputs, and the
// package turns those inputs into immutable version assignments for publishable
// modules.
package versioning
