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

// Details carries non-secret structured context for an infrastructure error.
//
// Values stored here may be rendered in diagnostics and logs. Implementations
// MUST NOT put raw secrets, tokens, or credentials into Details.
//
// Keep keys short and stable because callers may use them in tests, reports, or
// machine-readable diagnostics. Values should describe context, not duplicate
// the human Message.
type Details map[string]string

// Clone returns a detached copy of the details map.
//
// Empty maps normalize to nil so zero-value Error values stay compact and easy
// to compare in tests.
func (d Details) Clone() Details {
	if len(d) == 0 {
		return nil
	}
	out := make(Details, len(d))
	for key, value := range d {
		out[key] = value
	}
	return out
}
