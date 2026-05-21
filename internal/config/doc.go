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

// Package config loads ARCORIS Publisher manifest files from storage.
//
// The package owns file discovery, file reading, format detection, strict
// decoding, module-manifest path resolution, and construction of the effective
// resolved publication set. It deliberately does not build dependency graphs,
// assign versions, construct target repositories, invoke Git, run Go toolchain
// commands, or publish refs.
package config
