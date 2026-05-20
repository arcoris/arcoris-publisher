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

package remote

import (
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

func TestErrorCodesAreUniqueRemoteCodes(t *testing.T) {
	codes := []porterr.Code{
		CodeRepositoryNotFound,
		CodeAccessDenied,
		CodeAuthenticationFailed,
		CodeRateLimited,
		CodeBranchProtected,
		CodePullRequestFailed,
		CodeReleaseFailed,
	}
	seen := map[porterr.Code]bool{}

	for _, code := range codes {
		if code == "" {
			t.Fatalf("remote error code must not be empty")
		}
		if !strings.HasPrefix(code.String(), "remote_") {
			t.Fatalf("remote error code %q must use remote_ prefix", code)
		}
		if seen[code] {
			t.Fatalf("duplicate remote error code %q", code)
		}
		seen[code] = true
	}
}
