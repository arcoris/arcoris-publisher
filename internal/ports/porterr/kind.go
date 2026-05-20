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

package porterr

// Kind identifies the external infrastructure boundary that produced an error.
//
// Kind is intentionally broader than Code. It answers "which boundary failed?"
// while Code answers "which stable failure class happened inside that boundary?"
type Kind string

const (
	// KindProcess identifies errors raised while starting or waiting for an external process.
	KindProcess Kind = "process"
	// KindGit identifies errors raised by Git operations.
	KindGit Kind = "git"
	// KindFilesystem identifies errors raised by filesystem operations.
	KindFilesystem Kind = "filesystem"
	// KindGo identifies errors raised by the Go toolchain.
	KindGo Kind = "go"
	// KindRemote identifies errors raised by a remote hosting provider API.
	KindRemote Kind = "remote"
)

// String returns the stable string representation of the error kind.
func (k Kind) String() string {
	return string(k)
}
