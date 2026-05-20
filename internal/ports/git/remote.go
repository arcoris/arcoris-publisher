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

package git

import "context"

// RemoteClient describes Git remote transport operations such as clone, fetch,
// and push. Hosting-provider metadata APIs belong to the remote port package,
// not to this Git port.
type RemoteClient interface {
	// Clone copies remoteURL into dir according to opts.
	//
	// Implementations must redact opts.SensitiveValues from command rendering,
	// errors, and logs because remoteURL may contain credentials.
	Clone(ctx context.Context, remoteURL string, dir string, opts CloneOptions) error
	// Fetch updates refs from remote according to opts.
	//
	// RefSpecs are passed in order so callers can make deterministic fetch
	// requests when multiple refspecs overlap.
	Fetch(ctx context.Context, repoDir string, remote string, opts FetchOptions) error
	// Push updates remote refs using refspec and opts.
	//
	// Force and ForceWithLease are intentionally separate so workflow code can
	// prefer lease-protected publishing and reserve unguarded force pushes for
	// explicit operator decisions.
	Push(ctx context.Context, repoDir string, remote string, refspec RefSpec, opts PushOptions) error
}
