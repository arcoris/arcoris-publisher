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

import "testing"

func TestStrictDecoderDecodesStagingYAML(t *testing.T) {
	spec, err := (StrictDecoder{}).DecodeStaging([]byte(stagingYAML()), FormatYAML)
	if err != nil {
		t.Fatal(err)
	}
	if spec.APIVersion != "arcpub.arcoris.dev/v1alpha1" || len(spec.Modules) != 2 {
		t.Fatalf("unexpected staging spec: %#v", spec)
	}
}

func TestStrictDecoderDecodesModuleYAML(t *testing.T) {
	spec, err := (StrictDecoder{}).DecodeModule([]byte(controlModuleYAML()), FormatYAML)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Metadata.Name != "control" || len(spec.Publish.Entries) != 2 {
		t.Fatalf("unexpected module spec: %#v", spec)
	}
}

func TestStrictDecoderRejectsUnknownJSONFields(t *testing.T) {
	_, err := (StrictDecoder{}).DecodeStaging(
		stagingJSONWithUnknownField(),
		FormatJSON,
	)
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestYAMLDecoderRejectsDuplicateKeys(t *testing.T) {
	data := []byte("apiVersion: one\napiVersion: two\n")
	if _, err := (StrictDecoder{}).DecodeStaging(data, FormatYAML); err == nil {
		t.Fatal("expected duplicate key error")
	}
}

func TestYAMLDecoderRejectsUnsupportedInlineCollection(t *testing.T) {
	data := []byte(`apiVersion: arcpub.arcoris.dev/v1alpha1
kind: StagingManifest
metadata: {name: arcoris}
`)
	if _, err := (StrictDecoder{}).DecodeStaging(data, FormatYAML); err == nil {
		t.Fatal("expected unsupported inline collection error")
	}
}

func TestStrictDecoderRejectsUnsupportedFormat(t *testing.T) {
	if _, err := (StrictDecoder{}).DecodeStaging([]byte("{}"), Format("toml")); err == nil {
		t.Fatal("expected unsupported format error")
	}
}

func TestStrictDecoderRejectsTrailingJSON(t *testing.T) {
	_, err := (StrictDecoder{}).DecodeStaging(
		stagingJSONWithTrailingData(),
		FormatJSON,
	)
	if err == nil {
		t.Fatal("expected trailing JSON error")
	}
}
