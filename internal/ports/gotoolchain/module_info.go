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

// ModuleInfo describes the module metadata reported by the Go command.
//
// Replace points at the replacement module when a replace directive applies.
// It may itself contain only a subset of fields, matching Go's JSON output.
type ModuleInfo struct {
	// Path is the module path.
	Path string
	// Version is the selected module version.
	Version string
	// Dir is the module directory on disk when available.
	Dir string
	// Replace contains replacement module metadata when a replace directive applies.
	Replace *ModuleInfo
}

// HasReplace reports whether module replacement metadata is present.
func (m ModuleInfo) HasReplace() bool {
	return m.Replace != nil
}
