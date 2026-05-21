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

// PublishMode describes how target repository content is constructed.
type PublishMode string

// VersionPolicy describes how versions are assigned to modules.
type VersionPolicy string

// PushPolicy describes which Git push behavior is allowed.
type PushPolicy string

const (
	// PublishModeExplicitProjection means target content is built only from explicit module publish entries.
	PublishModeExplicitProjection PublishMode = "explicit-projection"

	// VersionPolicyReleaseTrain assigns one version to all publishable modules.
	VersionPolicyReleaseTrain VersionPolicy = "release-train"
	// VersionPolicySnapshot assigns snapshot or pseudo versions for non-release publication.
	VersionPolicySnapshot VersionPolicy = "snapshot"

	// PushPolicyFastForwardOnly forbids non-fast-forward target updates.
	PushPolicyFastForwardOnly PushPolicy = "fast-forward-only"
	// PushPolicyCreateOnly allows only creation of missing remote refs.
	PushPolicyCreateOnly PushPolicy = "create-only"
	// PushPolicyForceWithLease allows force-with-lease updates when explicitly supported by workflow.
	PushPolicyForceWithLease PushPolicy = "force-with-lease"
)

// ParsePublishMode validates a publication construction mode.
func ParsePublishMode(value string) (PublishMode, error) {
	switch PublishMode(value) {
	case PublishModeExplicitProjection:
		return PublishMode(value), nil
	default:
		return "", fmt.Errorf("unsupported publish mode %q", value)
	}
}

// ParseVersionPolicy validates a version assignment policy.
func ParseVersionPolicy(value string) (VersionPolicy, error) {
	switch VersionPolicy(value) {
	case VersionPolicyReleaseTrain, VersionPolicySnapshot:
		return VersionPolicy(value), nil
	default:
		return "", fmt.Errorf("unsupported versionPolicy %q", value)
	}
}

// ParsePushPolicy validates a push policy.
func ParsePushPolicy(value string) (PushPolicy, error) {
	switch PushPolicy(value) {
	case PushPolicyFastForwardOnly, PushPolicyCreateOnly, PushPolicyForceWithLease:
		return PushPolicy(value), nil
	default:
		return "", fmt.Errorf("unsupported pushPolicy %q", value)
	}
}
