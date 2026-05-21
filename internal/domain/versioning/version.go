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

// Version is a validated Go module tag or pseudo-version.
//
// Version intentionally uses the canonical module form with a leading "v", for
// example v0.1.0, v1.2.3-rc.1, or v0.0.0-20260521100000-abcdef123456. The type
// stores a string because version values cross adapter boundaries as tag names,
// go.mod requirements, and CLI arguments.
type Version string
