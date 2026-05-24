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

// Package report renders stable human-readable and machine-readable summaries
// for publisher plans and workflow results.
//
// The package transforms already-built internal models into output DTOs and
// text/JSON renderings. It does not load manifests, build plans, inspect Git,
// access the filesystem, execute workflow stages, rewrite module files, or
// publish repositories.
//
// Reports avoid local absolute filesystem paths by default because output may
// be uploaded to CI logs, attached to issues, or committed as artifacts.
package report
