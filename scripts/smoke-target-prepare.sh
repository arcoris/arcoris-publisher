#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"
LAB_ROOT="${TMPDIR:-/tmp}/arcpub-smoke-target-prepare"

echo "smoke-target-prepare: go test ./test/e2e -run TargetPrepare"
go test ./test/e2e -run TargetPrepare -count=1 -v

echo "smoke-target-prepare: lab target prepare"
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/setup.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/run-target-prepare.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/cleanup.sh

echo "smoke-target-prepare: ok"
