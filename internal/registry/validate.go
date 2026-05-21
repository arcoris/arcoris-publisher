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

package registry

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

// validateModules rejects impossible duplicate indexes before Registry is
// exposed to downstream packages.
func validateModules(modules []resolved.PublicationModule) error {
	collector := newIssueCollector()

	checkDuplicates(
		&collector,
		IssueDuplicateModuleName,
		"modules[%d].name",
		modules,
		func(module resolved.PublicationModule) string {
			return module.Name().String()
		},
	)
	checkDuplicates(
		&collector,
		IssueDuplicateModulePath,
		"modules[%d].module.path",
		modules,
		func(module resolved.PublicationModule) string {
			return module.ModulePath().String()
		},
	)
	checkDuplicates(
		&collector,
		IssueDuplicateSourceDir,
		"modules[%d].sourceDir",
		modules,
		func(module resolved.PublicationModule) string {
			return module.SourceDir().String()
		},
	)
	checkDuplicates(
		&collector,
		IssueDuplicateRepository,
		"modules[%d].repository",
		modules,
		func(module resolved.PublicationModule) string {
			return module.Repository().String()
		},
	)

	return collector.Err()
}

func checkDuplicates(
	collector *issueCollector,
	code IssueCode,
	pathFormat string,
	modules []resolved.PublicationModule,
	value func(resolved.PublicationModule) string,
) {
	seen := map[string]int{}
	for i, module := range modules {
		current := value(module)
		if first, ok := seen[current]; ok {
			collector.Add(
				code,
				fmt.Sprintf(pathFormat, i),
				"duplicates modules[%d] value %q",
				first,
				current,
			)
			continue
		}

		seen[current] = i
	}
}

// invalidPublicationSet reports a registry input that is structurally unusable.
func invalidPublicationSet(path string, format string, args ...any) error {
	collector := newIssueCollector()
	collector.Add(IssueInvalidPublicationSet, path, format, args...)
	return collector.Err()
}
