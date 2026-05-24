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

package porttest

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

func TestGoToolchainModTidyRunsHook(t *testing.T) {
	called := false
	fake := GoToolchain{
		ModTidyHook: func(context.Context, string) error {
			called = true
			return nil
		},
	}

	_, err := fake.ModTidy(context.Background(), "/module", gotoolchain.ModTidyOptions{})

	if err != nil {
		t.Fatalf("ModTidy() error = %v", err)
	}
	if !called {
		t.Fatal("hook was not called")
	}
}
