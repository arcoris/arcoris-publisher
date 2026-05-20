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

// Succeeded reports whether the result exit code is accepted by the provided list.
//
// It delegates to IsAllowedExitCode so the default empty-list policy stays
// identical for process specs and completed process results.
func (r Result) Succeeded(allowed []int) bool {
	return IsAllowedExitCode(r.ExitCode, allowed)
}
