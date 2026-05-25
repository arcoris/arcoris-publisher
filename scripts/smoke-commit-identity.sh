#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

echo "smoke-commit-identity: e2e commit identity scenarios"
go test ./test/e2e -run 'CommitIdentity|IdentityMissing' -count=1 -v

echo "smoke-commit-identity: lab missing identity scenario"
LAB_ROOT="${TMPDIR:-/tmp}/arcpub-smoke-commit-identity"
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/run-missing-identity.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/cleanup.sh

echo "smoke-commit-identity: ok"
