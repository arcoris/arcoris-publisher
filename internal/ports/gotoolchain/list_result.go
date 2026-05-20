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

// ListResult contains go list output and optional parsed package data.
//
// Stdout and Stderr preserve the command's raw streams. Packages is populated
// only when the adapter parsed JSON output, usually because ListOptions.JSON was
// set by the caller.
type ListResult struct {
	// Stdout contains raw standard output.
	Stdout []byte
	// Stderr contains raw standard error.
	Stderr []byte
	// Packages contains parsed package data when JSON output was requested.
	Packages []Package
}

// HasPackages reports whether parsed packages were attached to the result.
func (r ListResult) HasPackages() bool {
	return len(r.Packages) > 0
}
