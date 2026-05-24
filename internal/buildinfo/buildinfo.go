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

package buildinfo

const (
	defaultVersion = "dev"
	defaultCommit  = "unknown"
	defaultDate    = "unknown"
	defaultDirty   = "unknown"
)

// Version is the publisher version embedded by release builds.
//
// Build tooling may override this variable with -ldflags. Empty values are
// normalized to "dev" so local and test builds remain stable.
var Version = defaultVersion

// Commit is the source revision embedded by release builds.
//
// Build tooling may override this variable with -ldflags. Empty values are
// normalized to "unknown" because build metadata is optional for local builds.
var Commit = defaultCommit

// Date is the build timestamp embedded by release builds.
//
// Build tooling may override this variable with -ldflags. Empty values are
// normalized to "unknown" to avoid leaking machine-local fallback values.
var Date = defaultDate

// Dirty records whether the publisher source tree was dirty at build time.
//
// Build tooling may override this variable with -ldflags. Empty values are
// normalized to "unknown" because this package does not inspect Git itself.
var Dirty = defaultDirty

// Info is a normalized immutable snapshot of publisher build metadata.
//
// The fields are intentionally unexported so callers use normalized accessors
// and cannot accidentally distinguish unset values from defaults.
type Info struct {
	version string
	commit  string
	date    string
	dirty   string
}

// Current returns build metadata from the package-level ldflags variables.
func Current() Info {
	return Info{
		version: normalize(Version, defaultVersion),
		commit:  normalize(Commit, defaultCommit),
		date:    normalize(Date, defaultDate),
		dirty:   normalize(Dirty, defaultDirty),
	}
}

// Version returns the normalized publisher version.
func (i Info) Version() string { return normalize(i.version, defaultVersion) }

// Commit returns the normalized publisher source revision.
func (i Info) Commit() string { return normalize(i.commit, defaultCommit) }

// Date returns the normalized build timestamp.
func (i Info) Date() string { return normalize(i.date, defaultDate) }

// Dirty returns the normalized build dirty-state marker.
func (i Info) Dirty() string { return normalize(i.dirty, defaultDirty) }

// IsDev reports whether the publisher is a development build.
func (i Info) IsDev() bool { return i.Version() == defaultVersion }

// Map returns detached normalized metadata keyed by stable field names.
func (i Info) Map() map[string]string {
	return map[string]string{
		"version": i.Version(),
		"commit":  i.Commit(),
		"date":    i.Date(),
		"dirty":   i.Dirty(),
	}
}

func normalize(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
