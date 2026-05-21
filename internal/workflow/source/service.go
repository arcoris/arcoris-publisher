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

package source

import (
	"context"
	"fmt"
)

// Service inspects source repository state for publication plans.
type Service struct {
	// deps contains infrastructure ports required by source inspection.
	deps Dependencies

	// opts contains behavior toggles resolved at construction.
	opts Options
}

// New creates a source inspection service.
func New(deps Dependencies, opts Options) Service {
	if opts == (Options{}) {
		opts = DefaultOptions()
	}
	return Service{deps: deps, opts: opts}
}

// Inspect validates source repository state and returns a source snapshot.
func (s Service) Inspect(ctx context.Context, req Request) (Snapshot, error) {
	ins := inspector{deps: s.deps, opts: s.opts, request: req}
	snap, err := ins.inspect(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// validateDependencies rejects incomplete service wiring before request
// validation tries to use a nil port.
func (s Service) validateDependencies() error {
	if s.deps.Git == nil {
		return fmt.Errorf("source git dependency is required")
	}
	if s.deps.FS == nil {
		return fmt.Errorf("source filesystem dependency is required")
	}
	return nil
}
