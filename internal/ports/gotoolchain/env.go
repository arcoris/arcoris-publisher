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
