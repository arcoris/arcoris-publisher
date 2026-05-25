#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

echo "==> transaction e2e"
go test ./test/e2e -run 'TestPublishTransaction|TestPublishCandidate|TestPublishPromotion|TestPublishTag|TestPendingTransaction|TestRollbackCommand|TestRollbackFailure|TestTransactionsList' -count=1 -v

echo "transaction smoke passed"
