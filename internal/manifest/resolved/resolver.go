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

// resolver owns one resolution pass and its trace.
type resolver struct {
	input ResolveInput
	trace ResolutionTrace
}

// resolve binds staging modules to module manifests and validates the result.
func (r *resolver) resolve() (ResolveResult, error) {
	moduleByName, err := r.indexModuleManifests()
	if err != nil {
		return ResolveResult{}, err
	}

	resolvedModules, err := r.resolveModules(moduleByName)
	if err != nil {
		return ResolveResult{}, err
	}

	set := r.publicationSet(resolvedModules)
	if err := validatePublicationSet(set); err != nil {
		return ResolveResult{}, err
	}

	return ResolveResult{Set: set, Trace: r.trace}, nil
}
