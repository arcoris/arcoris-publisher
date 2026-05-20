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

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// duplicateModuleNameIssue describes a name claimed by two modules.
func duplicateModuleNameIssue(module manifest.ModuleName, other manifest.ModuleName) Issue {
	return Issue{
		Code:    IssueDuplicateModuleName,
		Module:  module,
		Other:   other,
		Message: fmt.Sprintf("duplicate module name %q already used by module %q", module, other),
	}
}

// duplicateModulePathIssue describes a public module path claimed twice.
func duplicateModulePathIssue(module manifest.ModuleName, other manifest.ModuleName, modulePath manifest.ModulePath) Issue {
	return Issue{
		Code:       IssueDuplicateModulePath,
		Module:     module,
		Other:      other,
		ModulePath: modulePath,
		Message:    fmt.Sprintf("duplicate module path %q already used by module %q", modulePath, other),
	}
}

// duplicateSourceDirIssue describes a staged source directory claimed twice.
func duplicateSourceDirIssue(module manifest.ModuleName, other manifest.ModuleName, sourceDir manifest.SourceDir) Issue {
	return Issue{
		Code:      IssueDuplicateSourceDir,
		Module:    module,
		Other:     other,
		SourceDir: sourceDir,
		Message:   fmt.Sprintf("duplicate source directory %q already used by module %q", sourceDir, other),
	}
}

// duplicateRepositoryIssue describes a target repository claimed twice.
func duplicateRepositoryIssue(module manifest.ModuleName, other manifest.ModuleName, repository manifest.RepositoryRef) Issue {
	return Issue{
		Code:       IssueDuplicateRepository,
		Module:     module,
		Other:      other,
		Repository: repository,
		Message:    fmt.Sprintf("duplicate repository %q already used by module %q", repository, other),
	}
}

// duplicateBranchMappingIssue describes a repeated source branch in one module.
func duplicateBranchMappingIssue(module manifest.ModuleName, branch manifest.BranchName) Issue {
	return Issue{
		Code:    IssueDuplicateBranchMapping,
		Module:  module,
		Branch:  branch,
		Message: fmt.Sprintf("module %q declares duplicate source branch mapping %q", module, branch),
	}
}
