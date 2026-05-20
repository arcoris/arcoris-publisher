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

package gotoolchain

import "time"

// EnvOptions configures go env execution.
type EnvOptions struct {
	// GoBinary overrides the executable name or path. Empty means "go".
	GoBinary string
	// Env contains additional KEY=VALUE environment assignments.
	Env []string
	// Timeout limits the tool invocation when greater than zero.
	Timeout time.Duration
}

// EnvResult contains selected Go environment values.
//
// Values is intentionally a map because different adapters may query different
// variables. Use HasValue when an empty string is a meaningful reported value.
type EnvResult struct {
	// Values maps Go environment variable names to their reported values.
	Values map[string]string
}

// Value returns a Go environment value by name or an empty string when missing.
func (r EnvResult) Value(name string) string {
	if r.Values == nil {
		return ""
	}
	return r.Values[name]
}

// HasValue reports whether a Go environment value was reported by name.
func (r EnvResult) HasValue(name string) bool {
	if r.Values == nil {
		return false
	}
	_, ok := r.Values[name]
	return ok
}
