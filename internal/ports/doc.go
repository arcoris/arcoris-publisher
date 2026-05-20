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

// Package ports contains infrastructure boundary contracts used by
// arcoris-publisher workflows.
//
// Ports describe capabilities supplied by the outside world: a filesystem, Git,
// the Go toolchain, process execution, and remote hosting providers. They are
// deliberately small and workflow-oriented so publisher logic can depend on
// behavior instead of concrete command-line tools or SDK clients.
//
// Packages under ports must not contain ARCORIS publishing business rules. They
// define vocabulary, data transfer objects, and error codes shared by workflow
// code and infrastructure adapters. Concrete implementations belong under
// internal/adapters and should translate implementation-specific failures into
// porterr.Error values with stable Kind and Code fields.
package ports
