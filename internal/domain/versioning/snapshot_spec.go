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

import "time"

// SnapshotSpec contains the source data required to create a deterministic
// development snapshot version.
type SnapshotSpec struct {
	// Base is the pseudo-version base. Empty defaults to v0.0.0.
	Base string
	// Time is the source commit timestamp used in the pseudo-version.
	Time time.Time
	// Commit is the source commit hash. At least 12 hexadecimal characters are required.
	Commit string
}
