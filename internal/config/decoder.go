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
	"bytes"
	"encoding/json"
	"fmt"

	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// Decoder converts serialized manifest bytes into syntax-level manifest specs.
type Decoder interface {
	DecodeStaging(data []byte, format Format) (staging.Spec, error)
	DecodeModule(data []byte, format Format) (modulemanifest.Spec, error)
}

// StrictDecoder rejects unknown fields and supports JSON plus the YAML subset
// used by ARCORIS Publisher manifests.
type StrictDecoder struct{}

// DecodeStaging decodes a top-level arcpub manifest spec.
func (d StrictDecoder) DecodeStaging(data []byte, format Format) (staging.Spec, error) {
	var spec staging.Spec
	if err := d.decode(data, format, &spec); err != nil {
		return staging.Spec{}, err
	}
	return spec, nil
}

// DecodeModule decodes a module-level arcpub.module manifest spec.
func (d StrictDecoder) DecodeModule(data []byte, format Format) (modulemanifest.Spec, error) {
	var spec modulemanifest.Spec
	if err := d.decode(data, format, &spec); err != nil {
		return modulemanifest.Spec{}, err
	}
	return spec, nil
}

// decode routes the selected manifest format to the strict decoder for that
// syntax. YAML is converted to JSON first so unknown-field checks stay shared.
func (StrictDecoder) decode(data []byte, format Format, out any) error {
	switch format {
	case FormatJSON:
		return decodeStrictJSON(data, out)
	case FormatYAML:
		jsonData, err := yamlToJSON(data)
		if err != nil {
			return err
		}
		return decodeStrictJSON(jsonData, out)
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

// decodeStrictJSON rejects unknown object fields and trailing JSON tokens before
// returning decoded manifest specs to domain validation.
func decodeStrictJSON(data []byte, out any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err == nil {
		return fmt.Errorf("unexpected trailing JSON data")
	}
	return nil
}
