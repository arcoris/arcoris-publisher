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

package gotoolchain

// Package describes one package returned by go list -json.
//
// The struct mirrors the subset of Go command JSON needed by publisher
// workflows. It intentionally omits many go list fields so adapters can remain
// stable across Go versions while still carrying import graph and module data.
type Package struct {
	// ImportPath is the package import path reported by the Go command.
	ImportPath string
	// Module is the module that provides the package.
	Module ModuleInfo
	// Imports are direct non-test package imports.
	Imports []string
	// TestImports are direct imports used by package tests.
	TestImports []string
	// Deps are transitive package dependencies when requested.
	Deps []string
	// Error describes a package loading error when go list reports one.
	Error *PackageError
}

// ModuleInfo describes the module metadata reported by the Go command.
//
// Replace points at the replacement module when a replace directive applies.
// It may itself contain only a subset of fields, matching Go's JSON output.
type ModuleInfo struct {
	// Path is the module path.
	Path string
	// Version is the selected module version.
	Version string
	// Dir is the module directory on disk when available.
	Dir string
	// Replace contains replacement module metadata when a replace directive applies.
	Replace *ModuleInfo
}

// PackageError describes a package loading error reported by the Go command.
type PackageError struct {
	// Err is the human-readable load error reported by go list.
	Err string
}

// HasReplace reports whether module replacement metadata is present.
func (m ModuleInfo) HasReplace() bool {
	return m.Replace != nil
}
