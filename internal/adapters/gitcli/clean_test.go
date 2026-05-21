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

package gitcli

import (
	"context"
	"testing"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestCleanFlags(t *testing.T) {
	tests := []struct {
		name string
		opts gitport.CleanOptions
		want []string
	}{
		{name: "noop", opts: gitport.CleanOptions{}, want: nil},
		{name: "untracked", opts: gitport.CleanOptions{RemoveUntracked: true}, want: []string{"clean", "-f"}},
		{name: "ignored only", opts: gitport.CleanOptions{RemoveIgnored: true, Directories: true}, want: []string{"clean", "-fdX"}},
		{
			name: "all forced",
			opts: gitport.CleanOptions{
				RemoveUntracked: true,
				RemoveIgnored:   true,
				Directories:     true,
				Force:           true,
			},
			want: []string{"clean", "-ffdx"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &fakeRunner{}
			client := New(runner, Options{})
			if err := client.Clean(context.Background(), "/repo", tt.opts); err != nil {
				t.Fatalf("Clean() error = %v", err)
			}
			if tt.want == nil {
				if len(runner.specs) != 0 {
					t.Fatalf("Clean() should not run command: %#v", runner.specs)
				}
				return
			}
			assertStringSlice(t, runner.specs[0].Args, tt.want)
		})
	}
}
