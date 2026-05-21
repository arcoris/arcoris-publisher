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

// Base returns the pseudo-version base.
func (s Snapshot) Base() Version { return s.base }

// Time returns the normalized UTC source timestamp.
func (s Snapshot) Time() time.Time { return s.time }

// Commit returns the normalized source commit hash.
func (s Snapshot) Commit() string { return s.commit }
