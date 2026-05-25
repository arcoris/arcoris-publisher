#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

run_step "resetting lab for happy path"
bash "$LAB_DIR/setup.sh"

report="$(lab_reports)/publish.json"
stderr="$(lab_logs)/publish.stderr"

run_step "running transactional publish"
code="$(run_common_arcpub_json "$report" "$stderr" publish)"
assert_code "$code" 0 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"
assert_transaction_status "$report" committed

tx_id="$(transaction_id_from_json "$report")"
[[ -n "$tx_id" ]] || fail "transaction id missing from publish report"

for repo in arcoris/foundation arcoris/control; do
  assert_published_repo "$repo"
  assert_clean_worktree "$(repository_worktree "$repo")"
done

run_step "saving transaction list/show"
"$(lab_arcpub)" transactions list --state-dir "$(lab_state_dir)" --output json >"$(lab_reports)/transactions-list.json"
"$(lab_arcpub)" transactions show "$tx_id" --state-dir "$(lab_state_dir)" --output json >"$(lab_reports)/transaction-show.json"
assert_json_if_python "$(lab_reports)/transactions-list.json"
assert_json_if_python "$(lab_reports)/transaction-show.json"
assert_no_path_leak "$(lab_reports)/transactions-list.json"
assert_no_path_leak "$(lab_reports)/transaction-show.json"

log "transaction: $tx_id"
log "publish report: $report"
