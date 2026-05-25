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

package preflight

import "arcoris.dev/arcoris-publisher/internal/workflow/source"

// Options configures read-only publish readiness checks.
type Options struct {
	// Source configures source inspection.
	Source source.Options

	// RemoteName is the target remote used by publish.
	RemoteName string

	// StateDir contains publish transaction journals and the publish lock.
	StateDir string
}

// DefaultOptions returns conservative preflight defaults.
func DefaultOptions() Options {
	return Options{RemoteName: "origin"}
}
