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

// TestOptions configures go test execution.
type TestOptions struct {
	CommonOptions
	// Patterns are package patterns passed to go test.
	Patterns []string
	// Race enables the race detector.
	Race bool
	// Count controls test caching when greater than zero.
	Count int
	// Short enables testing.Short for package tests.
	Short bool
	// Run filters test names with Go's -run regexp.
	Run string
	// Verbose enables verbose test output.
	Verbose bool
}

// TestResult contains go test output.
type TestResult struct {
	// Stdout contains raw standard output.
	Stdout []byte
	// Stderr contains raw standard error.
	Stderr []byte
}
