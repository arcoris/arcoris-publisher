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

// New creates a structured infrastructure error.
//
// New leaves Details empty and Temporary false. Call WithDetails and
// WithTemporary to enrich the error while preserving copy-on-write behavior.
func New(kind Kind, code Code, message string, cause error) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}
