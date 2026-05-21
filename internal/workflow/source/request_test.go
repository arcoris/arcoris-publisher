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

import "testing"

func TestRequestStoresPlanAndRoots(t *testing.T) {
	p := standardPlan(t)
	req := Request{
		Plan:          p,
		RepositoryDir: "/repo",
		StagingDir:    "/repo/staging",
	}

	if req.Plan.Empty() {
		t.Fatal("Plan is empty")
	}
	if req.RepositoryDir != "/repo" {
		t.Fatalf("RepositoryDir = %q", req.RepositoryDir)
	}
	if req.StagingDir != "/repo/staging" {
		t.Fatalf("StagingDir = %q", req.StagingDir)
	}
}
