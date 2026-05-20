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

// Visibility controls whether a module participates in external publication.
type Visibility string

const (
	// VisibilityPublic marks a module as externally publishable.
	VisibilityPublic Visibility = "public"
	// VisibilityInternal marks a module as known but not externally published.
	VisibilityInternal Visibility = "internal"
	// VisibilityDisabled marks a module as temporarily disabled.
	VisibilityDisabled Visibility = "disabled"
)

// ParseVisibility validates a module visibility. Empty visibility defaults to public.
func ParseVisibility(value string) (Visibility, error) {
	if value == "" {
		return VisibilityPublic, nil
	}
	switch Visibility(value) {
	case VisibilityPublic, VisibilityInternal, VisibilityDisabled:
		return Visibility(value), nil
	default:
		return "", fmt.Errorf("unsupported visibility %q", value)
	}
}
