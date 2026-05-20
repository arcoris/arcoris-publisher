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
//
// The type does not validate hash length or algorithm. Adapters should return
// the exact object name produced by their Git implementation.
type CommitHash string

// BranchName is a Git branch name.
//
// Values are stored without refs/heads/ unless a method explicitly documents a
// fully qualified ref.
type BranchName string

// TagName is a Git tag name.
//
// Values are stored without refs/tags/ unless a method explicitly documents a
// fully qualified ref.
type TagName string

// RefSpec is a Git refspec such as refs/heads/main:refs/heads/main.
//
// RefSpec is deliberately opaque because Git accepts several valid forms,
// including delete and force update syntax.
type RefSpec string
