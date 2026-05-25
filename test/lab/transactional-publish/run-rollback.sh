#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
tx_id="${1:-}"
[[ -n "$tx_id" ]] || fail "usage: $0 <transaction-id>"

first="$(lab_reports)/rollback.json"
second="$(lab_reports)/rollback-again.json"
err1="$(lab_logs)/rollback.stderr"
err2="$(lab_logs)/rollback-again.stderr"

run_step "rolling back $tx_id"
code="$(capture_arcpub "$first" "$err1" rollback --transaction "$tx_id" --state-dir "$(lab_state_dir)" --output json)"
assert_code "$code" 0 "$err1"
assert_json_if_python "$first"
assert_no_path_leak "$first"

run_step "rolling back $tx_id again to demonstrate idempotency"
code="$(capture_arcpub "$second" "$err2" rollback --transaction "$tx_id" --state-dir "$(lab_state_dir)" --output json)"
assert_code "$code" 0 "$err2"
assert_json_if_python "$second"
assert_no_path_leak "$second"
log "rollback reports: $first, $second"
