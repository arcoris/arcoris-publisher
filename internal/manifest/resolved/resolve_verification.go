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
	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
)

// resolveVerification applies built-ins, staging defaults, then module overrides.
func (r *resolver) resolveVerification(
	path string,
	mm modulemanifest.Manifest,
) manifest.VerificationPolicy {
	tracePath := path + ".verification"

	policy := manifest.BuiltInVerificationPolicy()
	r.trace.AddBuiltInDefault(tracePath, "built-in", "verification")

	policy = manifest.MergeVerification(policy, r.input.Staging.Defaults().Verification())
	r.trace.AddStagingDefault(
		tracePath,
		"staging defaults applied",
		"defaults.verification",
	)

	policy = manifest.MergeVerification(policy, mm.Verification())
	r.trace.AddModuleManifest(
		tracePath,
		"module overrides applied",
		string(mm.Metadata().Name())+".verification",
	)

	return policy
}
