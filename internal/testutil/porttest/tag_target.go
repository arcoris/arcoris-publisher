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

package porttest

import (
	"context"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

// TagTargetHash returns the fake repository's stable commit for an existing
// tag. The in-memory Git model represents created refs with CommitHash, so the
// same object identity is used for ownership checks during rollback tests.
func (g *Git) TagTargetHash(_ context.Context, _ string, tag git.TagName) (git.CommitHash, bool, error) {
	if g.TagExistsError != nil {
		return "", false, g.TagExistsError
	}
	if !g.Tags[tag] {
		return "", false, nil
	}
	return g.CommitHash, true, nil
}
