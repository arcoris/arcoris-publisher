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

// APIVersion identifies the manifest schema group and version.
type APIVersion string

// Kind identifies the concrete manifest resource kind.
type Kind string

const (
	// APIVersionV1Alpha1 is the first unstable ARCORIS Publisher manifest API.
	APIVersionV1Alpha1 APIVersion = "arcpub.arcoris.dev/v1alpha1"

	// KindStagingManifest identifies the top-level arcpub.yaml resource.
	KindStagingManifest Kind = "StagingManifest"
	// KindModuleManifest identifies the module-level arcpub.module.yaml resource.
	KindModuleManifest Kind = "ModuleManifest"
)

// ValidateAPIVersion verifies that value is a supported manifest API version.
func ValidateAPIVersion(value string) (APIVersion, error) {
	if value == "" {
		return "", NewFieldError(IssueMissingField, "apiVersion", "apiVersion is required")
	}
	version := APIVersion(value)
	if version != APIVersionV1Alpha1 {
		return "", NewFieldError(IssueInvalidAPIVersion, "apiVersion", "unsupported apiVersion %q", value)
	}
	return version, nil
}

// ValidateKind verifies that value matches the expected manifest kind.
func ValidateKind(value string, expected Kind) (Kind, error) {
	if value == "" {
		return "", NewFieldError(IssueMissingField, "kind", "kind is required")
	}
	kind := Kind(value)
	if kind != expected {
		return "", NewFieldError(IssueInvalidKind, "kind", "expected kind %q, got %q", expected, value)
	}
	return kind, nil
}
