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

// assignmentsValidator validates assignment state without mutating the receiver.
type assignmentsValidator struct {
	assignments Assignments
	indexes     assignmentIndexes
	issues      []Issue
}

// newAssignmentsValidator prepares validation for one assignment set.
func newAssignmentsValidator(assignments Assignments) assignmentsValidator {
	return assignmentsValidator{
		assignments: assignments,
		indexes:     newAssignmentIndexes(len(assignments.items)),
	}
}

// validate checks policy, item validity, and duplicate lookup keys.
func (v *assignmentsValidator) validate() (assignmentIndexes, error) {
	v.validatePolicy()
	for index, item := range v.assignments.items {
		v.validateItem(index, item)
	}
	if len(v.issues) > 0 {
		return assignmentIndexes{}, &ValidationError{Issues: v.issues}
	}
	return v.indexes, nil
}
