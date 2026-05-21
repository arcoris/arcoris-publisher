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
	"strings"
	"testing"
)

func TestServiceRejectsMissingDependencies(t *testing.T) {
	_, err := New(Dependencies{}, DefaultOptions()).Inspect(
		context.Background(),
		standardRequest(t),
	)
	if err == nil || !strings.Contains(err.Error(), "git dependency") {
		t.Fatalf("expected missing git dependency error, got %v", err)
	}

	_, err = New(Dependencies{Git: cleanGit()}, DefaultOptions()).Inspect(
		context.Background(),
		standardRequest(t),
	)
	if err == nil || !strings.Contains(err.Error(), "filesystem dependency") {
		t.Fatalf("expected missing filesystem dependency error, got %v", err)
	}
}
