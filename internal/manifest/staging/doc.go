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

// Package staging defines the top-level arcpub.yaml manifest resource.
//
// The staging manifest owns publication topology and global publication policy:
// source repository, staging root, module list, target repositories, branch
// routing, version/push/tag/provenance policy, and defaults. It deliberately
// does not describe module-local published files or directories; those belong to
// package module and arcpub.module.yaml.
package staging
