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

// IsAllowedExitCode reports whether code is accepted by the allowed set.
//
// If allowed is empty, only code 0 is accepted. This default matches the normal
// operating-system process convention while still allowing callers to explicitly
// accept commands where a non-zero exit code is a meaningful result.
func IsAllowedExitCode(code int, allowed []int) bool {
	if len(allowed) == 0 {
		return code == 0
	}
	for _, item := range allowed {
		if item == code {
			return true
		}
	}
	return false
}
