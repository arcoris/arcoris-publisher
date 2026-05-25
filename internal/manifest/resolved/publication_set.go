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

package resolved

import "arcoris.dev/arcoris-publisher/internal/manifest"

// PublicationSet is the effective resolved publication model.
type PublicationSet struct {
	metadata manifest.Metadata
	source   manifest.Source
	target   manifest.TargetPolicy
	publish  manifest.PublishPolicy
	modules  []PublicationModule
}

// Metadata returns the publication set metadata.
func (s PublicationSet) Metadata() manifest.Metadata { return s.metadata }

// Source returns the effective source repository declaration.
func (s PublicationSet) Source() manifest.Source { return s.source }

// Target returns the effective target preparation policy.
func (s PublicationSet) Target() manifest.TargetPolicy { return s.target }

// Publish returns the effective global publication policy.
func (s PublicationSet) Publish() manifest.PublishPolicy { return s.publish }

// Modules returns detached effective publication modules.
func (s PublicationSet) Modules() []PublicationModule {
	out := make([]PublicationModule, len(s.modules))
	copy(out, s.modules)
	return out
}
