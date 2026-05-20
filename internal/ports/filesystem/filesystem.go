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

// Package filesystem defines the infrastructure port for safe filesystem tree
// operations required by publisher workflows.
//
// The port is intentionally higher level than os or io/fs. Publisher workflows
// operate on whole trees, safety roots, include/exclude filters, and deterministic
// hashes; adapters hide platform-specific path behavior and filesystem error
// details behind this contract.
//
// Implementations should treat every path argument as an external boundary:
// normalize it consistently, reject operations that escape configured safety
// roots, honor context cancellation where possible, and return structured
// porterr.Error values using this package's error codes.
package filesystem

// FileSystem groups filesystem capabilities required by publisher workflows.
type FileSystem interface {
	Reader
	Writer
	Cleaner
	Copier
	Hasher
}
