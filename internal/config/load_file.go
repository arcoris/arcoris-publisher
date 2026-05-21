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

import "context"

type manifestFile struct {
	Path   string
	Format Format
	Data   []byte
}

// readManifestFile performs the common file-loading pipeline shared by staging
// and module manifests: normalize the path, read bytes, then detect the syntax
// format from the final path.
func (l Loader) readManifestFile(
	ctx context.Context,
	path string,
) (manifestFile, error) {
	path, err := normalizeInputPath(path)
	if err != nil {
		return manifestFile{}, err
	}

	data, err := l.reader.ReadFile(ctx, path)
	if err != nil {
		return manifestFile{}, readError(path, err)
	}

	format, err := DetectFormat(path)
	if err != nil {
		return manifestFile{}, err
	}

	return manifestFile{
		Path:   path,
		Format: format,
		Data:   data,
	}, nil
}
