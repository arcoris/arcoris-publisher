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

package registry

// Validate verifies registry invariants after construction.
//
// A Registry returned by New has already been validated. Validate is primarily
// useful for tests, zero-value checks, and defensive checks on manually
// assembled registries inside this package. It first replays public
// construction rules from the module list, then compares every internal index
// against the expected index set.
func (r Registry) Validate() error {
	expected, err := New(r.Modules())
	if err != nil {
		return err
	}
	return newIndexValidator(r, expected).validate()
}
