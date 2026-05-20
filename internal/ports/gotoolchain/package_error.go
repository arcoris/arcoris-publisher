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

package gotoolchain

// PackageError describes a package loading error reported by the Go command.
//
// The Go command may return package objects with embedded load errors even when
// it emits useful partial graph data. Keeping the error as data lets higher
// layers decide whether partial information is acceptable for their workflow.
type PackageError struct {
	// Err is the human-readable load error reported by go list.
	Err string
}
