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

// TagClient describes Git tag operations required by release publishing.
type TagClient interface {
	TagExists(ctx context.Context, repoDir string, tag TagName) (bool, error)
	CreateTag(ctx context.Context, repoDir string, tag TagName, target CommitHash, opts TagOptions) error
	PushTag(ctx context.Context, repoDir string, remote string, tag TagName, opts PushOptions) error
}
