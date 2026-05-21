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

package publish

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"testing"
)

func TestResultPublishedAndDetachedTags(t *testing.T) {
	r := Result{modules: []ModuleResult{{
		module: manifest.ModuleName("control"),
		pushed: true,
		tags:   []git.TagName{"v0.1.0"},
	}}}
	if !r.Published() {
		t.Fatalf("expected published result")
	}
	mods := r.Modules()
	tags := mods[0].Tags()
	tags[0] = "mutated"
	if r.Modules()[0].Tags()[0] != "v0.1.0" {
		t.Fatalf("tags were not detached")
	}
}
