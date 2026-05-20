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

// ListOptions configures go list execution.
type ListOptions struct {
	CommonOptions
	// Patterns are package patterns passed to go list.
	Patterns []string
	// Deps includes transitive dependencies in the result.
	Deps bool
	// Test asks go list to include test variants.
	Test bool
	// JSON asks the adapter to parse go list -json output into Packages.
	JSON bool
}
