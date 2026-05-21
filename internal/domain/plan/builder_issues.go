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

package plan

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// missingModuleIssue reports publish order output that the registry cannot resolve.
func missingModuleIssue(module manifest.ModuleName) Issue {
	return Issue{
		Code:    IssueMissingModule,
		Module:  module,
		Message: fmt.Sprintf("module %q is missing from registry", module),
	}
}

// missingVersionIssue reports a publishable module without an assignment.
func missingVersionIssue(module manifest.ModuleName) Issue {
	return Issue{
		Code:    IssueMissingVersion,
		Module:  module,
		Message: fmt.Sprintf("module %q has no assigned version", module),
	}
}

// invalidBranchIssue reports a module whose branch mappings cannot become a plan.
func invalidBranchIssue(module manifest.ModuleName, err error) Issue {
	return Issue{
		Code:    IssueInvalidBranch,
		Module:  module,
		Message: fmt.Sprintf("module %q has invalid branch plan: %v", module, err),
	}
}

// duplicateModuleIssue reports two plans with the same module name.
func duplicateModuleIssue(module manifest.ModuleName, previous int) Issue {
	return Issue{
		Code:    IssueDuplicateModule,
		Module:  module,
		Message: fmt.Sprintf("duplicate module plan %q already exists at index %d", module, previous),
	}
}

// duplicateModulePathIssue reports two plans with the same module path.
func duplicateModulePathIssue(module manifest.ModuleName, path manifest.ModulePath, previous int) Issue {
	return Issue{
		Code:    IssueDuplicateModulePath,
		Module:  module,
		Path:    path,
		Message: fmt.Sprintf("duplicate module path %q already exists at index %d", path, previous),
	}
}

// duplicateRepositoryIssue reports two plans with the same target repository.
func duplicateRepositoryIssue(module manifest.ModuleName, repository manifest.RepositoryRef, previous int) Issue {
	return Issue{
		Code:       IssueDuplicateRepository,
		Module:     module,
		Repository: repository,
		Message:    fmt.Sprintf("duplicate repository %q already exists at index %d", repository, previous),
	}
}
