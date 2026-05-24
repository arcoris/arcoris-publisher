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

// Package app contains high-level publisher use cases.
//
// The application layer coordinates existing packages only: it loads manifests,
// builds indexes and topology, assigns versions, creates a plan, and delegates
// execution to the workflow runner. Filesystem copying, Git publication, Go
// module rewriting, and verification remain in workflow stages and ports.
package app
