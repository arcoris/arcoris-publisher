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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

type builder struct {
	modules []manifest.Module

	registry Registry
	issues   []Issue
}

func newBuilder(modules []manifest.Module) builder {
	return builder{modules: cloneModules(modules)}
}

func (b builder) build() (Registry, error) {
	b.registry = newEmptyRegistry(b.modules)
	for index, module := range b.modules {
		b.indexModule(index, module)
	}

	if len(b.issues) > 0 {
		return Registry{}, &ValidationError{Issues: b.issues}
	}
	return b.registry, nil
}
