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

package versioning

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// version resolves the single version assigned by the selected policy.
func (b assignmentBuilder) version(policy manifest.VersionPolicy) (Version, error) {
	switch policy {
	case manifest.VersionPolicyReleaseTrain:
		return b.releaseTrainVersion()
	case manifest.VersionPolicySnapshot:
		return b.snapshotVersion()
	default:
		return "", validationErrorf(IssueUnsupportedPolicy, "policy", "unsupported version policy %q", policy)
	}
}

// releaseTrainVersion parses the release version from the assignment spec.
func (b assignmentBuilder) releaseTrainVersion() (Version, error) {
	version, err := ParseVersion(b.spec.Release)
	if err != nil {
		return "", validationErrorf(IssueInvalidReleaseVersion, "release", "invalid release version: %v", err)
	}
	return version, nil
}

// snapshotVersion builds a pseudo-version from the snapshot spec.
func (b assignmentBuilder) snapshotVersion() (Version, error) {
	snapshot, err := NewSnapshot(b.spec.Snapshot)
	if err != nil {
		return "", validationErrorf(IssueInvalidSnapshot, "snapshot", "invalid snapshot: %v", err)
	}
	return snapshot.Version(), nil
}
