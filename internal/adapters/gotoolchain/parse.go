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

import (
	"bytes"
	"encoding/json"
	"io"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// jsonPackage mirrors the subset of one go list -json package object we need.
//
// The Go command exposes many more fields. Keeping this DTO intentionally small
// prevents adapter code from depending on unstable or currently unused output.
type jsonPackage struct {
	ImportPath  string
	Module      *jsonModule
	Imports     []string
	TestImports []string
	Deps        []string
	Error       *goport.PackageError
}

// jsonModule mirrors the recursive module object emitted by go list -json.
//
// Replace uses the same JSON shape as a normal module, so the DTO is recursive
// and convertModule handles that recursion explicitly.
type jsonModule struct {
	Path    string
	Version string
	Dir     string
	Replace *jsonModule
}

// parsePackages decodes the stream produced by go list -json.
//
// go list emits one JSON object per package rather than a single array, so this
// function repeatedly decodes objects until EOF and returns the first malformed
// JSON error to the caller.
func parsePackages(data []byte) ([]goport.Package, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	packages := []goport.Package{}
	for {
		var item jsonPackage
		if err := dec.Decode(&item); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		packages = append(packages, convertPackage(item))
	}
	return packages, nil
}

// convertPackage converts one JSON package object into the port package model.
//
// Slices are detached because the returned result belongs to callers, not the
// decoder buffer or temporary DTO values.
func convertPackage(item jsonPackage) goport.Package {
	return goport.Package{
		ImportPath:  item.ImportPath,
		Module:      convertModule(item.Module),
		Imports:     append([]string(nil), item.Imports...),
		TestImports: append([]string(nil), item.TestImports...),
		Deps:        append([]string(nil), item.Deps...),
		Error:       item.Error,
	}
}

// convertModule converts the recursive JSON module shape into the port type.
func convertModule(m *jsonModule) goport.ModuleInfo {
	if m == nil {
		return goport.ModuleInfo{}
	}
	out := goport.ModuleInfo{Path: m.Path, Version: m.Version, Dir: m.Dir}
	if m.Replace != nil {
		repl := convertModule(m.Replace)
		out.Replace = &repl
	}
	return out
}
