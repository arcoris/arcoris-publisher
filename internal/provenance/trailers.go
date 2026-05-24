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

package provenance

import "strings"

// Trailer is one deterministic Git commit trailer.
type Trailer struct {
	Key   string
	Value string
}

// Trailers is an ordered collection of Git commit trailers.
type Trailers []Trailer

// BuildTrailers builds commit trailers from the same source as file
// provenance. The order is stable so commit messages remain reproducible.
func BuildTrailers(input Input) Trailers {
	entries := EntriesFromSourceModule(input.SourceModule)
	build := input.Build

	return Trailers{
		{Key: "Arcoris-Source-Repository", Value: input.Plan.Source().Repository().String()},
		{Key: "Arcoris-Source-Commit", Value: string(input.Source.Repository().Head())},
		{Key: "Arcoris-Source-Branch", Value: string(input.Source.Repository().Branch())},
		{Key: "Arcoris-Module", Value: input.Module.Name().String()},
		{Key: "Arcoris-Module-Path", Value: input.Module.ModulePath().String()},
		{Key: "Arcoris-Version", Value: input.Module.Version().String()},
		{Key: "Arcoris-Target-Repository", Value: input.Module.Repository().String()},
		{Key: "Arcoris-Target-Branches", Value: strings.Join(input.targetBranches(), ",")},
		{Key: "Arcoris-Publish-Mode", Value: string(input.Plan.PublishPolicy().Mode())},
		{Key: "Arcoris-Push-Policy", Value: string(input.Plan.PublishPolicy().PushPolicy())},
		{Key: "Arcoris-Tag-Policy", Value: string(input.Plan.PublishPolicy().Tags().Mode())},
		{Key: "Arcoris-Publisher-Version", Value: build.Version()},
		{Key: "Arcoris-Source-Dir", Value: input.Module.SourceDir().String()},
		{Key: "Arcoris-Source-Hash", Value: input.SourceModule.Hash().String()},
		{Key: "Arcoris-Projection-Hash", Value: ProjectionHash(entries)},
	}
}

// Render renders trailers in order with one key/value pair per line.
func (t Trailers) Render() string {
	var b strings.Builder
	for _, trailer := range t {
		b.WriteString(sanitizeTrailerKey(trailer.Key))
		b.WriteString(": ")
		b.WriteString(sanitizeTrailerValue(trailer.Value))
		b.WriteByte('\n')
	}
	return b.String()
}

func sanitizeTrailerKey(key string) string {
	return strings.TrimSpace(replaceLineBreaks(key))
}

func sanitizeTrailerValue(value string) string {
	return strings.TrimSpace(replaceLineBreaks(value))
}

func replaceLineBreaks(value string) string {
	replacer := strings.NewReplacer(
		"\r\n", " ",
		"\r", " ",
		"\n", " ",
	)
	return replacer.Replace(value)
}
