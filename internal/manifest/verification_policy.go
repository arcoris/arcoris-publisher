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

// VerificationPolicy is a complete resolved verification policy.
type VerificationPolicy struct {
	localReplacePolicy LocalReplacePolicy
	goPolicy           GoVerificationPolicy
}

// GoVerificationPolicy is a complete resolved Go verification policy.
type GoVerificationPolicy struct {
	workspaceMode GoWorkspaceMode
	list          bool
	test          bool
	tidy          bool
	patterns      []string
}

// BuiltInVerificationPolicy returns the safe built-in verification defaults.
func BuiltInVerificationPolicy() VerificationPolicy {
	return VerificationPolicy{
		localReplacePolicy: LocalReplacePolicyForbid,
		goPolicy: GoVerificationPolicy{
			workspaceMode: GoWorkspaceModeOff,
			list:          true,
			test:          true,
			tidy:          true,
			patterns:      []string{"./..."},
		},
	}
}

// LocalReplacePolicy returns the resolved local replace policy.
func (p VerificationPolicy) LocalReplacePolicy() LocalReplacePolicy { return p.localReplacePolicy }

// Go returns the resolved Go verification policy.
func (p VerificationPolicy) Go() GoVerificationPolicy { return p.goPolicy }

// WorkspaceMode returns the resolved Go workspace mode.
func (p GoVerificationPolicy) WorkspaceMode() GoWorkspaceMode { return p.workspaceMode }

// List reports whether go list verification is enabled.
func (p GoVerificationPolicy) List() bool { return p.list }

// Test reports whether go test verification is enabled.
func (p GoVerificationPolicy) Test() bool { return p.test }

// Tidy reports whether go mod tidy stability verification is enabled.
func (p GoVerificationPolicy) Tidy() bool { return p.tidy }

// Patterns returns detached Go package patterns for verification.
func (p GoVerificationPolicy) Patterns() []string { return cloneStrings(p.patterns) }
