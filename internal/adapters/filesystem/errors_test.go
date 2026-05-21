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

package filesystem

import (
	"errors"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

func TestFSErrorAndPathErrorUseFilesystemKind(t *testing.T) {
	cause := errors.New("failed")
	err := pathError(fsport.CodePathNotFound, "missing", cause, "/repo")

	var perr *porterr.Error
	if !errors.As(err, &perr) {
		t.Fatalf("pathError() type = %T", err)
	}
	if perr.Kind != porterr.KindFilesystem || perr.Code != fsport.CodePathNotFound || perr.Cause != cause {
		t.Fatalf("pathError() = %#v", perr)
	}
	if perr.Details["path"] != "/repo" || !isPortError(err) {
		t.Fatalf("pathError() details/isPortError mismatch: %#v", perr)
	}
}
