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

package remote

import "testing"

func TestBranchProtectionBlocksDirectPush(t *testing.T) {
	tests := []struct {
		name       string
		protection BranchProtection
		want       bool
	}{
		{name: "unprotected", protection: BranchProtection{}, want: false},
		{name: "protected only", protection: BranchProtection{Protected: true}, want: false},
		{name: "pull request only", protection: BranchProtection{RequiresPullRequest: true}, want: false},
		{name: "protected pull request", protection: BranchProtection{Protected: true, RequiresPullRequest: true}, want: true},
		{name: "status checks only", protection: BranchProtection{Protected: true, RequiresStatusChecks: true}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.protection.BlocksDirectPush(); got != tt.want {
				t.Fatalf("BlocksDirectPush() = %v, want %v", got, tt.want)
			}
		})
	}
}
