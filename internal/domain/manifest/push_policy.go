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

// PushPolicy describes how target branches and tags may be pushed.
type PushPolicy string

const (
	// PushPolicyFastForwardOnly allows normal non-force pushes only.
	PushPolicyFastForwardOnly PushPolicy = "fast-forward-only"
	// PushPolicyCreateOnly allows creating new remote refs but not updating existing refs.
	PushPolicyCreateOnly PushPolicy = "create-only"
	// PushPolicyForceWithLease allows force-with-lease pushes when explicitly requested.
	PushPolicyForceWithLease PushPolicy = "force-with-lease"
)

// ParsePushPolicy validates a push policy.
func ParsePushPolicy(value string) (PushPolicy, error) {
	switch PushPolicy(value) {
	case PushPolicyFastForwardOnly, PushPolicyCreateOnly, PushPolicyForceWithLease:
		return PushPolicy(value), nil
	default:
		return "", fmt.Errorf("unsupported push policy %q", value)
	}
}
