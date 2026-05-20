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

// WithDetails returns a copy of the error with detached details attached.
//
// The receiver is not mutated. This allows adapters to define base errors and
// attach operation-specific context without sharing maps between errors.
func (e *Error) WithDetails(details Details) *Error {
	if e == nil {
		return nil
	}
	copy := *e
	copy.Details = details.Clone()
	return &copy
}
