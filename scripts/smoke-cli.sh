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

echo "smoke-cli: go test ./..."
go test ./...

echo "smoke-cli: go vet ./..."
go vet ./...

echo "smoke-cli: build arcpub"
go build -o "$tmp/arcpub" ./cmd/arcpub

echo "smoke-cli: arcpub help"
"$tmp/arcpub" help >/dev/null

echo "smoke-cli: arcpub version"
"$tmp/arcpub" version >/dev/null

echo "smoke-cli: arcpub version --output json"
"$tmp/arcpub" version --output json >"$tmp/version.json"
if command -v python3 >/dev/null 2>&1; then
	python3 -m json.tool "$tmp/version.json" >/dev/null
else
	echo "python3 not found; skipping JSON validation"
fi

echo "smoke-cli: arcpub completion bash"
"$tmp/arcpub" completion bash >"$tmp/completion.bash"
if [[ ! -s "$tmp/completion.bash" ]]; then
	echo "bash completion output is empty" >&2
	exit 1
fi

echo "smoke-cli: ok"
