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

package process

// WithSensitiveValues returns a copy of the spec with detached sensitive values.
//
// The values are copied so later caller mutations cannot change redaction
// behavior for a spec already passed through workflow code.
func (s Spec) WithSensitiveValues(values ...string) Spec {
	copy := s
	copy.SensitiveValues = append([]string(nil), values...)
	return copy
}

// WithAllowedExitCodes returns a copy of the spec with detached success codes.
//
// Passing no codes resets the spec to the default success policy: only exit code
// 0 is accepted by IsAllowedExitCode.
func (s Spec) WithAllowedExitCodes(codes ...int) Spec {
	copy := s
	copy.AllowedExitCodes = append([]int(nil), codes...)
	return copy
}
