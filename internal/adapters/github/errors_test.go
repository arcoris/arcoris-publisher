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

package github

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

func TestRemoteErrorAndCodeHelpers(t *testing.T) {
	err := remoteError(remoteport.CodeAccessDenied, "denied", errors.New("cause"), porterr.Details{"path": "/repo"})
	code, ok := remoteCode(err)

	if !ok || code != remoteport.CodeAccessDenied || !isRemoteCode(err, remoteport.CodeAccessDenied) {
		t.Fatalf("remoteCode()/isRemoteCode() = %s, %v", code, ok)
	}
	if _, ok := remoteCode(errors.New("plain")); ok {
		t.Fatalf("remoteCode() should reject plain errors")
	}
}

func TestWrapRemoteOperationErrorPreservesSpecificRemoteCodes(t *testing.T) {
	cause := remoteError(remoteport.CodeRepositoryNotFound, "missing", nil, nil)
	if got := wrapRemoteOperationError(remoteport.CodeReleaseFailed, "release failed", cause, nil); got != cause {
		t.Fatalf("wrapRemoteOperationError() should preserve repository_not_found")
	}

	err := wrapRemoteOperationError(remoteport.CodeReleaseFailed, "release failed", errors.New("plain"), nil)
	assertPortCode(t, err, remoteport.CodeReleaseFailed)
}
