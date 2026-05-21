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

package staging

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// Validate checks aggregate-level arcpub.yaml invariants.
func (m Manifest) Validate() error {
	var collector manifest.IssueCollector
	names := make(map[manifest.ModuleName]int, len(m.modules))
	sourceDirs := make(map[manifest.SourceDir]int, len(m.modules))
	for i, mod := range m.modules {
		path := fmt.Sprintf("modules[%d]", i)
		if prev, ok := names[mod.Name()]; ok {
			collector.Add(manifest.IssueDuplicateValue, path+".name", "duplicate module name %q previously declared at modules[%d]", mod.Name(), prev)
		} else {
			names[mod.Name()] = i
		}
		if prev, ok := sourceDirs[mod.SourceDir()]; ok {
			collector.Add(manifest.IssueDuplicateValue, path+".sourceDir", "duplicate sourceDir %q previously declared at modules[%d]", mod.SourceDir(), prev)
		} else {
			sourceDirs[mod.SourceDir()] = i
		}
	}
	return collector.Err()
}
