#!/usr/bin/env bash
# Copyright 2026 The ARCORIS Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

go test ./...
go vet ./...
go build -o "$tmp/arcpub" ./cmd/arcpub

"$tmp/arcpub" help >/dev/null
"$tmp/arcpub" version >/dev/null
"$tmp/arcpub" version --output json >"$tmp/version.json"
if command -v python3 >/dev/null 2>&1; then
	python3 -m json.tool "$tmp/version.json" >/dev/null
else
	echo "python3 not found; skipping JSON validation"
fi

"$tmp/arcpub" completion bash >"$tmp/completion.bash"
test -s "$tmp/completion.bash"

echo "smoke-cli: ok"
