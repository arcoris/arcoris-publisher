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

package runtime

// Runtime is the production composition root for the publisher binary.
//
// Runtime stores normalized options and dependencies and exposes builders for
// application and CLI layers. It is intentionally small: command routing remains
// in internal/cli, and publication execution remains in internal/workflow.
type Runtime struct {
	opts Options
	deps Dependencies
}

// New creates a Runtime with production adapters derived from opts.
func New(opts Options) Runtime {
	opts = normalizeOptions(opts)
	return Runtime{opts: opts, deps: NewDependencies(opts)}
}

// NewWithDependencies creates a Runtime using caller-supplied dependencies and
// filling any omitted collaborators with production defaults.
func NewWithDependencies(deps Dependencies, opts Options) Runtime {
	opts = normalizeOptions(opts)
	return Runtime{opts: opts, deps: normalizeDependencies(deps, opts)}
}

// Options returns the normalized runtime options.
func (r Runtime) Options() Options {
	opts := r.opts
	opts.Env = copyStrings(opts.Env)
	return opts
}

// Dependencies returns the normalized runtime dependencies.
func (r Runtime) Dependencies() Dependencies { return r.deps }
