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

// WorkspaceMode controls how Go workspace mode is configured for toolchain calls.
//
// The zero value is intentionally invalid. Callers should choose
// WorkspaceDefault when they want the Go command's normal discovery behavior or
// WorkspaceOff when publishing must be isolated from ambient go.work files.
type WorkspaceMode string

const (
	// WorkspaceDefault leaves workspace selection to the Go command.
	WorkspaceDefault WorkspaceMode = "default"
	// WorkspaceOff forces single-module mode through GOWORK=off.
	WorkspaceOff WorkspaceMode = "off"
)

// String returns the stable string representation of the workspace mode.
func (m WorkspaceMode) String() string {
	return string(m)
}

// Valid reports whether the mode is one of the supported workspace modes.
func (m WorkspaceMode) Valid() bool {
	switch m {
	case WorkspaceDefault, WorkspaceOff:
		return true
	default:
		return false
	}
}
