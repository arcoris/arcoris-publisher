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

// ModTidyOptions configures go mod tidy execution.
type ModTidyOptions struct {
	CommonOptions
	// Compat is the Go version passed through -compat when non-empty.
	Compat string
}

// ModTidyResult contains go mod tidy output.
type ModTidyResult struct {
	// Stdout contains raw standard output.
	Stdout []byte
	// Stderr contains raw standard error.
	Stderr []byte
}
