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

package target

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestCheckCommitIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values map[string]string
		errs   map[string]error
		want   string
	}{
		{
			name: "configured",
			values: map[string]string{
				"user.name":  "ARCORIS Test",
				"user.email": "arcoris-test@example.invalid",
			},
		},
		{name: "missing both", want: CommitIdentityCodeMissingBoth},
		{
			name:   "missing name",
			values: map[string]string{"user.email": "arcoris-test@example.invalid"},
			want:   CommitIdentityCodeMissingName,
		},
		{
			name:   "missing email",
			values: map[string]string{"user.name": "ARCORIS Test"},
			want:   CommitIdentityCodeMissingEmail,
		},
		{
			name: "read error",
			errs: map[string]error{"user.name": errors.New("config failed")},
			want: CommitIdentityCodeReadFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fakeGit := porttest.NewGit()
			for key, value := range tt.values {
				fakeGit.ConfigValues[key] = value
			}
			for key, err := range tt.errs {
				fakeGit.ConfigErrors[key] = err
			}

			check := CheckCommitIdentity(context.Background(), fakeGit, "/repo")
			if got := check.Code(); got != tt.want {
				t.Fatalf("Code() = %q, want %q", got, tt.want)
			}
			if tt.want == "" && !check.Passed() {
				t.Fatal("Passed() = false")
			}
			if tt.want != "" && check.Passed() {
				t.Fatal("Passed() = true")
			}
		})
	}
}
