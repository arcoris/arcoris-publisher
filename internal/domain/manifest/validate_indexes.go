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

import "fmt"

// recordModuleName records module-name uniqueness.
func (v *manifestValidator) recordModuleName(path string, name ModuleName, index int) {
	if previous, exists := v.moduleNames[name]; exists {
		v.issues = append(v.issues, duplicateIssue(IssueDuplicateModule, path+".name", "module name", name, previous))
		return
	}
	v.moduleNames[name] = index
}

// recordModulePath records public Go module path uniqueness.
func (v *manifestValidator) recordModulePath(path string, modulePath ModulePath, index int) {
	if previous, exists := v.modulePaths[modulePath]; exists {
		v.issues = append(v.issues, duplicateIssue(IssueDuplicatePath, path+".module_path", "module path", modulePath, previous))
		return
	}
	v.modulePaths[modulePath] = index
}

// recordSourceDir records staged source directory uniqueness.
func (v *manifestValidator) recordSourceDir(path string, sourceDir SourceDir, index int) {
	if previous, exists := v.sourceDirs[sourceDir]; exists {
		v.issues = append(v.issues, duplicateIssue(IssueDuplicatePath, path+".source_dir", "source directory", sourceDir, previous))
		return
	}
	v.sourceDirs[sourceDir] = index
}

// recordRepository records target repository uniqueness.
func (v *manifestValidator) recordRepository(path string, repository RepositoryRef, index int) {
	if previous, exists := v.repositories[repository]; exists {
		v.issues = append(v.issues, duplicateIssue(IssueDuplicatePath, path+".repository", "target repository", repository, previous))
		return
	}
	v.repositories[repository] = index
}

// duplicateIssue creates a standard duplicate aggregate validation issue.
func duplicateIssue[T ~string](code IssueCode, path string, label string, value T, previous int) Issue {
	return Issue{
		Code:    code,
		Path:    path,
		Message: fmt.Sprintf("duplicate %s %q already used by modules[%d]", label, value, previous),
	}
}
