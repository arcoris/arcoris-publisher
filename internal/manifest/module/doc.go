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

// Package module defines the module-level arcpub.module.yaml manifest resource.
//
// A module manifest owns the local publication contract of one module: module
// identity, internal dependencies, explicit file/directory publication entries,
// and local verification overrides. It deliberately does not decide target
// repository, branch routing, push policy, version policy, or tag policy; those
// belong to the top-level staging manifest.
package module
