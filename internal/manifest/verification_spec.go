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

type text = string
type ptr = *text
type g = GoVerifySpec

// VerificationSpec is a partial verification-policy declaration.
//
// VerificationSpec is intentionally an override shape. Missing fields mean
// "inherit from a lower-precedence default", not false or empty.
type VerificationSpec struct {
	LocalReplacePolicy ptr `json:"localReplacePolicy,omitempty" yaml:"localReplacePolicy,omitempty"`
	Go                 g   `json:"go,omitempty" yaml:"go,omitempty"`
}

// GoVerifySpec is the partial Go verification declaration.
type GoVerifySpec struct {
	WorkspaceMode *string  `json:"workspaceMode,omitempty" yaml:"workspaceMode,omitempty"`
	List          *bool    `json:"list,omitempty" yaml:"list,omitempty"`
	Test          *bool    `json:"test,omitempty" yaml:"test,omitempty"`
	Tidy          *bool    `json:"tidy,omitempty" yaml:"tidy,omitempty"`
	Patterns      []string `json:"patterns,omitempty" yaml:"patterns,omitempty"`
}
