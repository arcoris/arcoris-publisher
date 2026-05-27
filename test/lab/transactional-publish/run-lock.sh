#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
require_cmd git

show_report="$(lab_reports)/transactions-lock-show.json"
show_stderr="$(lab_logs)/transactions-lock-show.stderr"
clear_report="$(lab_reports)/transactions-lock-clear.json"
clear_stderr="$(lab_logs)/transactions-lock-clear.stderr"
refused_report="$(lab_reports)/transactions-lock-clear-refused.json"
refused_stderr="$(lab_logs)/transactions-lock-clear-refused.stderr"

run_step "publishing once to create a committed transaction journal"
bash "$LAB_DIR/run-happy-path.sh"
tx_id="$(transaction_id_from_json "$(lab_reports)/publish.json")"
[[ -n "$tx_id" ]] || fail "could not read transaction id from publish report"
tx_path="$(lab_state_dir)/transactions/$tx_id.json"
assert_file_exists "$tx_path"

run_step "creating an operator-inspection lock for the committed transaction"
cat >"$(lab_state_dir)/publish.lock" <<LOCK
transaction=$tx_id
pid=1
startedAt=2026-01-01T00:00:00Z
command=publish
LOCK

run_step "showing current publish lock"
code="$(capture_arcpub "$show_report" "$show_stderr" \
  transactions lock show \
  --state-dir "$(lab_state_dir)" \
  --output json)"
assert_code "$code" 0 "$show_stderr"
assert_json_if_python "$show_report"
assert_no_path_leak "$show_report"
assert_contains "$show_report" '"status": "present"'

run_step "clearing committed transaction lock with explicit confirmation"
code="$(capture_arcpub "$clear_report" "$clear_stderr" \
  transactions lock clear \
  --state-dir "$(lab_state_dir)" \
  --transaction "$tx_id" \
  --confirm "$tx_id" \
  --output json)"
assert_code "$code" 0 "$clear_stderr"
assert_json_if_python "$clear_report"
assert_no_path_leak "$clear_report"
assert_contains "$clear_report" '"status": "cleared"'
assert_contains "$clear_report" '"reason": "cleared"'
assert_contains "$clear_report" '"postClearState": "ready_for_publish"'
[[ ! -e "$(lab_state_dir)/publish.lock" ]] || fail "publish lock was not cleared"
assert_file_exists "$tx_path"

run_step "creating active pending transaction lock that clear must refuse"
cat >"$(lab_state_dir)/transactions/tx-pending-lock-lab.json" <<'JSON'
{
  "schemaVersion": 1,
  "id": "tx-pending-lock-lab",
  "status": "pending",
  "version": "v0.1.0",
  "startedAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "modules": []
}
JSON
cat >"$(lab_state_dir)/publish.lock" <<'LOCK'
transaction=tx-pending-lock-lab
pid=1
startedAt=2026-01-01T00:00:00Z
command=publish
LOCK

code="$(capture_arcpub "$refused_report" "$refused_stderr" \
  transactions lock clear \
  --state-dir "$(lab_state_dir)" \
  --transaction tx-pending-lock-lab \
  --confirm tx-pending-lock-lab \
  --output json)"
assert_code "$code" 1 "$refused_stderr"
assert_json_if_python "$refused_report"
assert_no_path_leak "$refused_report"
assert_contains "$refused_report" '"status": "refused"'
assert_contains "$refused_report" '"reason": "active_transaction"'
assert_file_exists "$(lab_state_dir)/publish.lock"
assert_file_exists "$(lab_state_dir)/transactions/tx-pending-lock-lab.json"
rm -f "$(lab_state_dir)/publish.lock"

log "transactions lock show report: $show_report"
log "transactions lock clear report: $clear_report"
log "transactions lock refused report: $refused_report"
