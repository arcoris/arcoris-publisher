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

import (
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// ResolveInput contains validated staging and module manifests to bind.
type ResolveInput struct {
	Staging staging.Manifest
	Modules []modulemanifest.Manifest
}

// ResolveResult contains the effective publication set and value origin trace.
type ResolveResult struct {
	Set   PublicationSet
	Trace ResolutionTrace
}
