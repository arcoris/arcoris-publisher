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

	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// LoadStaging loads and validates one top-level arcpub.yaml/arcpub.json file.
func (l Loader) LoadStaging(
	ctx context.Context,
	path string,
) (staging.Manifest, error) {
	file, err := l.readManifestFile(ctx, path)
	if err != nil {
		return staging.Manifest{}, err
	}

	spec, err := l.decoder.DecodeStaging(file.Data, file.Format)
	if err != nil {
		return staging.Manifest{}, decodeError(file.Path, err, "staging")
	}

	manifest, err := staging.New(spec)
	if err != nil {
		return staging.Manifest{}, invalidManifestError(file.Path, err, "staging")
	}

	return manifest, nil
}
