#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

run_step "resetting lab for promotion failure"
bash "$LAB_DIR/setup.sh"
write_reject_hook "$(bare_repo arcoris/control)" promotion-control

report="$(lab_reports)/promotion-failure.json"
stderr="$(lab_logs)/promotion-failure.stderr"

run_step "running publish with final branch rejection"
code="$(run_common_arcpub_json "$report" "$stderr" publish)"
assert_code "$code" 1 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"
assert_transaction_status "$report" rolled_back
assert_no_final_refs
for repo in arcoris/foundation arcoris/control; do
  assert_clean_worktree "$(repository_worktree "$repo")"
done

tx_id="$(transaction_id_from_json "$report")"
[[ -n "$tx_id" ]] || fail "transaction id missing from failure report"
bash "$LAB_DIR/run-rollback.sh" "$tx_id"
log "promotion failure rollback demonstrated: $tx_id"
