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

package manifest

// cloneStrings returns a detached copy of a string slice.
func cloneStrings(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	return out
}

// CloneBranchMappings returns a detached copy of branch mappings.
func CloneBranchMappings(in []BranchMapping) []BranchMapping {
	out := make([]BranchMapping, len(in))
	copy(out, in)
	return out
}

// CloneModuleNames returns a detached copy of module names.
func CloneModuleNames(in []ModuleName) []ModuleName {
	out := make([]ModuleName, len(in))
	copy(out, in)
	return out
}

// ClonePublishEntries returns a detached copy of publish entries.
func ClonePublishEntries(in []PublishEntry) []PublishEntry {
	out := make([]PublishEntry, len(in))
	copy(out, in)
	return out
}
