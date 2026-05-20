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

package porterr

// Code is a stable machine-readable infrastructure error code.
//
// Concrete port packages define their own code constants. For example, the Git
// port defines push and ref-related codes, while the process port defines
// timeout and process startup codes.
//
// Codes are part of the workflow contract. Rename them only with a migration
// path because callers may branch on exact strings in retry, reporting, or
// operator guidance logic.
type Code string

// String returns the stable string representation of the error code.
func (c Code) String() string {
	return string(c)
}
