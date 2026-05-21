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

import (
	"errors"
	"fmt"
)

// validateItem checks one assignment item and indexes it when valid.
func (v *assignmentsValidator) validateItem(index int, item ModuleVersion) {
	path := fmt.Sprintf("items[%d]", index)
	if err := item.validate(); err != nil {
		v.addItemValidationError(path, err)
		return
	}
	v.indexItem(index, path, item)
}

// addItemValidationError rewrites nested item paths into assignment item paths.
func (v *assignmentsValidator) addItemValidationError(path string, err error) {
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		v.addIssue(Issue{Code: IssueInvalidAssignment, Path: path, Message: err.Error()})
		return
	}
	for _, issue := range validationErr.Issues {
		issue.Path = path + "." + issue.Path
		v.addIssue(issue)
	}
}

// indexItem adds one valid item to module and path lookup indexes.
func (v *assignmentsValidator) indexItem(index int, path string, item ModuleVersion) {
	if previous, exists := v.indexes.byModule[item.Module()]; exists {
		v.addIssue(Issue{
			Code:    IssueDuplicateModule,
			Path:    path + ".module",
			Message: fmt.Sprintf("duplicate module %q already assigned by items[%d]", item.Module(), previous),
		})
	} else {
		v.indexes.byModule[item.Module()] = index
	}

	if previous, exists := v.indexes.byPath[item.ModulePath()]; exists {
		v.addIssue(Issue{
			Code:    IssueDuplicateModulePath,
			Path:    path + ".module_path",
			Message: fmt.Sprintf("duplicate module path %q already assigned by items[%d]", item.ModulePath(), previous),
		})
	} else {
		v.indexes.byPath[item.ModulePath()] = index
	}
}
