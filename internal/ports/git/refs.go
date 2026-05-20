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

// CommitHash is a Git commit object hash.
type CommitHash string

// BranchName is a Git branch name.
type BranchName string

// TagName is a Git tag name.
type TagName string

// RefSpec is a Git refspec such as refs/heads/main:refs/heads/main.
type RefSpec string

// String returns the commit hash as a string.
func (h CommitHash) String() string { return string(h) }

// String returns the branch name as a string.
func (b BranchName) String() string { return string(b) }

// String returns the tag name as a string.
func (t TagName) String() string { return string(t) }

// String returns the refspec as a string.
func (r RefSpec) String() string { return string(r) }
