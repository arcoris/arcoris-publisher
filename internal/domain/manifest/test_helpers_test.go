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

func validSpec() Spec {
	return Spec{
		Version: "v1",
		Source:  SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"},
		Policy:  PolicySpec{VersionPolicy: "release-train", PushPolicy: "fast-forward-only"},
		Modules: []ModuleSpec{
			{
				Name:       "foundation",
				ModulePath: "arcoris.dev/foundation",
				SourceDir:  "staging/src/arcoris.dev/foundation",
				Repository: "arcoris/foundation",
				Branches:   []BranchMappingSpec{{Source: "main", Target: "main"}},
			},
			{
				Name:         "control",
				ModulePath:   "arcoris.dev/control",
				SourceDir:    "staging/src/arcoris.dev/control",
				Repository:   "arcoris/control",
				Branches:     []BranchMappingSpec{{Source: "main", Target: "main"}},
				Dependencies: []string{"foundation"},
			},
		},
	}
}
