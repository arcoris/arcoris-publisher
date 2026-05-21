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
	"context"

	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
)

// LoadModule loads and validates one arcpub.module.yaml/arcpub.module.json file.
func (l Loader) LoadModule(
	ctx context.Context,
	path string,
) (modulemanifest.Manifest, error) {
	file, err := l.readManifestFile(ctx, path)
	if err != nil {
		return modulemanifest.Manifest{}, err
	}

	spec, err := l.decoder.DecodeModule(file.Data, file.Format)
	if err != nil {
		return modulemanifest.Manifest{}, decodeError(file.Path, err, "module")
	}

	manifest, err := modulemanifest.New(spec)
	if err != nil {
		return modulemanifest.Manifest{}, invalidManifestError(file.Path, err, "module")
	}

	return manifest, nil
}
