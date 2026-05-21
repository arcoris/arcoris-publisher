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

package config

import (
	"path/filepath"
	"strings"
)

// Format identifies a supported serialized manifest format.
type Format string

const (
	// FormatYAML is used for .yaml and .yml manifest files.
	FormatYAML Format = "yaml"
	// FormatJSON is used for .json manifest files.
	FormatJSON Format = "json"
)

// DetectFormat returns the manifest format implied by path's extension.
func DetectFormat(path string) (Format, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		return FormatYAML, nil
	case ".json":
		return FormatJSON, nil
	default:
		return "", unsupportedFormatError(path)
	}
}

// StagingManifestNames returns the default top-level manifest names in lookup order.
func StagingManifestNames() []string {
	return []string{
		"arcpub.yaml",
		"arcpub.yml",
		"arcpub.json",
	}
}

// ModuleManifestNames returns conventional module-level manifest names in
// lookup order.
func ModuleManifestNames() []string {
	return []string{
		"arcpub.module.yaml",
		"arcpub.module.yml",
		"arcpub.module.json",
	}
}
