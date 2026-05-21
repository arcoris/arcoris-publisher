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

// Validate verifies plan-level invariants after construction.
//
// A Plan returned by New has already been validated. Validate is primarily
// useful for tests, zero-value checks, and defensive checks on manually
// assembled plans.
func (p Plan) Validate() error {
	validator := newPlanValidator(p)
	return validator.validate()
}
