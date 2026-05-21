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

package manifest

import "fmt"

// LocalReplacePolicy controls whether local replace directives are allowed.
type LocalReplacePolicy string

// GoWorkspaceMode controls Go workspace usage during verification.
type GoWorkspaceMode string

const (
	// LocalReplacePolicyForbid rejects local replace directives in published modules.
	LocalReplacePolicyForbid LocalReplacePolicy = "forbid"
	// LocalReplacePolicyWarn reports local replace directives as warnings.
	LocalReplacePolicyWarn LocalReplacePolicy = "warn"
	// LocalReplacePolicyAllow allows local replace directives.
	LocalReplacePolicyAllow LocalReplacePolicy = "allow"

	// GoWorkspaceModeOff verifies modules with GOWORK=off.
	GoWorkspaceModeOff GoWorkspaceMode = "off"
	// GoWorkspaceModeDefault lets the Go toolchain use its default workspace behavior.
	GoWorkspaceModeDefault GoWorkspaceMode = "default"
)

// ParseLocalReplacePolicy validates a local replace policy.
func ParseLocalReplacePolicy(value string) (LocalReplacePolicy, error) {
	switch LocalReplacePolicy(value) {
	case LocalReplacePolicyForbid, LocalReplacePolicyWarn, LocalReplacePolicyAllow:
		return LocalReplacePolicy(value), nil
	default:
		return "", fmt.Errorf("unsupported localReplacePolicy %q", value)
	}
}

// ParseGoWorkspaceMode validates a Go workspace mode.
func ParseGoWorkspaceMode(value string) (GoWorkspaceMode, error) {
	switch GoWorkspaceMode(value) {
	case GoWorkspaceModeOff, GoWorkspaceModeDefault:
		return GoWorkspaceMode(value), nil
	default:
		return "", fmt.Errorf("unsupported go.workspaceMode %q", value)
	}
}
