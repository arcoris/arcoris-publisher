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
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// resolveManifestPath applies per-module manifest path override before defaults.
func (r *resolver) resolveManifestPath(path string, sm staging.Module) manifest.RelativePath {
	if override, ok := sm.ManifestPathOverride(); ok {
		r.trace.Add(path+".manifest", override.String(), ValueSource{Kind: SourceStagingModule, Path: path + ".manifest"})
		return override
	}
	value := r.input.Staging.Defaults().ModuleManifest().Path()
	r.trace.Add(path+".manifest", value.String(), ValueSource{Kind: SourceStagingDefault, Path: "defaults.moduleManifest.path"})
	return value
}

// resolveVisibility applies per-module visibility override before the public default.
func (r *resolver) resolveVisibility(path string, sm staging.Module) manifest.Visibility {
	if override, ok := sm.VisibilityOverride(); ok {
		r.trace.Add(path+".visibility", override.String(), ValueSource{Kind: SourceStagingModule, Path: path + ".visibility"})
		return override
	}
	r.trace.Add(path+".visibility", manifest.VisibilityPublic.String(), ValueSource{Kind: SourceBuiltInDefault, Path: "visibility"})
	return manifest.VisibilityPublic
}

// resolveBranches applies module overrides, then defaults, then source default branch.
func (r *resolver) resolveBranches(path string, sm staging.Module) []manifest.BranchMapping {
	if override, ok := sm.BranchesOverride(); ok {
		r.trace.Add(path+".branches", fmt.Sprintf("%d mapping(s)", len(override)), ValueSource{Kind: SourceStagingModule, Path: path + ".branches"})
		return override
	}
	if r.input.Staging.Defaults().BranchesSet() {
		branches := r.input.Staging.Defaults().Branches()
		r.trace.Add(path+".branches", fmt.Sprintf("%d mapping(s)", len(branches)), ValueSource{Kind: SourceStagingDefault, Path: "defaults.branches"})
		return branches
	}
	branch := r.input.Staging.Source().DefaultBranch()
	mapping, _ := manifest.NewBranchMapping(manifest.BranchMappingSpec{Source: branch.String(), Target: branch.String()})
	r.trace.Add(path+".branches", fmt.Sprintf("%s->%s", branch, branch), ValueSource{Kind: SourceBuiltInDefault, Path: "source.defaultBranch"})
	return []manifest.BranchMapping{mapping}
}

// resolveVerification applies built-ins, staging defaults, then module overrides.
func (r *resolver) resolveVerification(path string, mm modulemanifest.Manifest) manifest.VerificationPolicy {
	policy := manifest.BuiltInVerificationPolicy()
	r.trace.Add(path+".verification", "built-in", ValueSource{Kind: SourceBuiltInDefault, Path: "verification"})
	policy = manifest.MergeVerification(policy, r.input.Staging.Defaults().Verification())
	r.trace.Add(path+".verification", "staging defaults applied", ValueSource{Kind: SourceStagingDefault, Path: "defaults.verification"})
	policy = manifest.MergeVerification(policy, mm.Verification())
	r.trace.Add(path+".verification", "module overrides applied", ValueSource{Kind: SourceModuleManifest, Path: string(mm.Metadata().Name()) + ".verification"})
	return policy
}
