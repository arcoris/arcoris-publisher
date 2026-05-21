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

package versioning

// validate checks that a single module assignment has all required keys.
func (v ModuleVersion) validate() error {
	var issues []Issue
	if v.module == "" {
		issues = append(issues, Issue{Code: IssueInvalidAssignment, Path: "module", Message: "module is required"})
	}
	if v.modulePath == "" {
		issues = append(issues, Issue{Code: IssueInvalidAssignment, Path: "module_path", Message: "module path is required"})
	}
	if v.version == "" {
		issues = append(issues, Issue{Code: IssueInvalidAssignment, Path: "version", Message: "version is required"})
	} else if _, err := ParseVersion(v.version.String()); err != nil {
		issues = append(issues, Issue{Code: IssueInvalidAssignment, Path: "version", Message: "version is invalid: " + err.Error()})
	}
	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
