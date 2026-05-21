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

package plan

// validateInputs checks each upstream domain projection before combining them.
func (b *builder) validateInputs() error {
	if err := b.registryValue.Validate(); err != nil {
		return validationErrorf(IssueInvalidRegistry, "", "invalid registry: %s", err)
	}
	if err := b.graphValue.Validate(); err != nil {
		return validationErrorf(IssueInvalidGraph, "", "invalid graph: %s", err)
	}
	if err := b.assignments.Validate(); err != nil {
		return validationErrorf(IssueMissingVersion, "", "invalid version assignments: %s", err)
	}
	return nil
}
