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

// CommonOptions contains options shared by Go toolchain operations.
type CommonOptions struct {
	// GoBinary overrides the executable name or path. Empty means "go".
	GoBinary string
	// WorkspaceMode controls GOWORK behavior for module-scoped commands.
	WorkspaceMode WorkspaceMode
	// Env contains additional KEY=VALUE environment assignments.
	Env []string
	// Tags are build tags passed to commands that compile packages.
	Tags []string
	// Timeout limits the tool invocation when greater than zero.
	Timeout time.Duration
	// PrivateModules configures GOPRIVATE-style module patterns.
	PrivateModules []string
	// Proxy configures GOPROXY when non-empty.
	Proxy string
	// SumDB configures GOSUMDB when non-empty.
	SumDB string
	// SensitiveValues are raw values that adapters must redact in diagnostics.
	SensitiveValues []string
}
