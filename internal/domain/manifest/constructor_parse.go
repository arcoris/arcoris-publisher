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

// parseManifestVersion validates the manifest schema version field.
func parseManifestVersion(spec Spec) (Version, error) {
	version, err := ParseVersion(spec.Version)
	if err != nil {
		return "", validationErrorf(IssueUnsupportedVersion, "version", "invalid manifest version: %v", err)
	}
	return version, nil
}

// parseManifestSource validates the authoritative source declaration.
func parseManifestSource(spec Spec) (Source, error) {
	source, err := NewSource(spec.Source)
	if err != nil {
		return Source{}, validationErrorf(IssueInvalidSource, "source", "invalid source: %v", err)
	}
	return source, nil
}

// parseManifestPolicy validates the global publication policy.
func parseManifestPolicy(spec Spec) (Policy, error) {
	policy, err := NewPolicy(spec.Policy)
	if err != nil {
		return Policy{}, validationErrorf(IssueInvalidPolicy, "policy", "invalid policy: %v", err)
	}
	return policy, nil
}

// parseManifestModules validates module declarations in order.
func parseManifestModules(specs []ModuleSpec) ([]Module, error) {
	modules := make([]Module, 0, len(specs))
	for i, spec := range specs {
		module, err := NewModule(spec)
		if err != nil {
			return nil, validationErrorf(IssueInvalidModule, issuePath("modules", i), "invalid module: %v", err)
		}
		modules = append(modules, module)
	}
	return modules, nil
}
