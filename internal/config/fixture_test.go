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

package config

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestFixtureMinimalLoadsPublicationSet(t *testing.T) {
	set, err := NewLoader(LoaderOptions{}).LoadPublicationSet(
		context.Background(),
		fixtureManifestPath("minimal"),
	)
	if err != nil {
		t.Fatal(err)
	}

	modules := set.Modules()
	if len(modules) != 2 {
		t.Fatalf("got %d modules", len(modules))
	}
	if modules[0].Name().String() != "foundation" {
		t.Fatalf("unexpected first module: %s", modules[0].Name())
	}
	if modules[1].Name().String() != "control" {
		t.Fatalf("unexpected second module: %s", modules[1].Name())
	}

	assertControlDependsOnFoundation(t, modules[1].Dependencies())
}

func TestFixtureInvalidTreesAreRejected(t *testing.T) {
	cases := []struct {
		name string
		code Code
	}{
		{name: "invalid-name-mismatch", code: CodeResolveFailed},
		{name: "unknown-dependency", code: CodeResolveFailed},
		{name: "missing-module-manifest", code: CodeReadFailed},
		{name: "unsafe-module-manifest-path", code: CodeInvalidManifest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewLoader(LoaderOptions{}).LoadPublicationSet(
				context.Background(),
				fixtureManifestPath(tc.name),
			)
			if err == nil {
				t.Fatal("expected fixture to be rejected")
			}

			var configErr *Error
			if !errors.As(err, &configErr) {
				t.Fatalf("expected config error, got %T", err)
			}
			if configErr.Code != tc.code {
				t.Fatalf("got code %s, want %s", configErr.Code, tc.code)
			}
		})
	}
}

func fixtureManifestPath(name string) string {
	return filepath.Join("testdata", name, "arcpub.yaml")
}
