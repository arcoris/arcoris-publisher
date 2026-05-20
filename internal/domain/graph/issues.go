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

package graph

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// duplicateModuleIssue reports a repeated module name in graph input.
func duplicateModuleIssue(name manifest.ModuleName, index int) Issue {
	return Issue{
		Code:    IssueDuplicateModule,
		Module:  name,
		Message: fmt.Sprintf("duplicate module %q at index %d", name, index),
	}
}

// selfDependencyIssue reports a direct module self edge.
func selfDependencyIssue(module manifest.ModuleName) Issue {
	return Issue{
		Code:    IssueSelfDependency,
		Module:  module,
		Message: fmt.Sprintf("module %q cannot depend on itself", module),
	}
}

// duplicateDependencyIssue reports a repeated direct dependency declaration.
func duplicateDependencyIssue(module manifest.ModuleName, dependency manifest.ModuleName) Issue {
	return Issue{
		Code:       IssueDuplicateDependency,
		Module:     module,
		Dependency: dependency,
		Message:    fmt.Sprintf("module %q declares duplicate dependency %q", module, dependency),
	}
}

// unknownDependencyIssue reports a dependency target missing from the graph.
func unknownDependencyIssue(module manifest.ModuleName, dependency manifest.ModuleName) Issue {
	return Issue{
		Code:       IssueUnknownDependency,
		Module:     module,
		Dependency: dependency,
		Message:    fmt.Sprintf("module %q depends on unknown module %q", module, dependency),
	}
}

// cycleIssue reports one dependency cycle.
func cycleIssue(cycle Cycle) Issue {
	return Issue{
		Code:    IssueCycle,
		Message: "dependency cycle detected: " + cycle.String(),
		Cycle:   cycle,
	}
}
