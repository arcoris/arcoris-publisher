#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
require_cmd git

dry_report="$(lab_reports)/transactions-prune-dry-run.json"
dry_stderr="$(lab_logs)/transactions-prune-dry-run.stderr"
prune_report="$(lab_reports)/transactions-prune.json"
prune_stderr="$(lab_logs)/transactions-prune.stderr"

run_step "publishing once to create a committed transaction journal"
bash "$LAB_DIR/run-happy-path.sh"
tx_id="$(transaction_id_from_json "$(lab_reports)/publish.json")"
[[ -n "$tx_id" ]] || fail "could not read transaction id from publish report"
tx_path="$(lab_state_dir)/transactions/$tx_id.json"
assert_file_exists "$tx_path"

run_step "creating rollback_failed sample journal that prune must preserve"
mkdir -p "$(lab_state_dir)/transactions"
cat >"$(lab_state_dir)/transactions/tx-rollback-failed-lab.json" <<'JSON'
{
  "schemaVersion": 1,
  "id": "tx-rollback-failed-lab",
  "status": "rollback_failed",
  "rollbackStatus": "failed",
  "version": "v0.1.0",
  "startedAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-01-01T00:00:00Z",
  "modules": []
}
JSON

run_step "previewing committed transaction prune"
code="$(capture_arcpub "$dry_report" "$dry_stderr" \
  transactions prune \
  --state-dir "$(lab_state_dir)" \
  --status committed \
  --dry-run \
  --output json)"
assert_code "$code" 0 "$dry_stderr"
assert_json_if_python "$dry_report"
assert_no_path_leak "$dry_report"
assert_contains "$dry_report" '"status": "dry_run"'
assert_file_exists "$tx_path"

run_step "pruning committed transaction journal"
code="$(capture_arcpub "$prune_report" "$prune_stderr" \
  transactions prune \
  --state-dir "$(lab_state_dir)" \
  --status committed \
  --output json)"
assert_code "$code" 0 "$prune_stderr"
assert_json_if_python "$prune_report"
assert_no_path_leak "$prune_report"
assert_contains "$prune_report" '"status": "completed"'
[[ ! -e "$tx_path" ]] || fail "committed transaction journal was not pruned: $tx_path"
assert_file_exists "$(lab_state_dir)/transactions/tx-rollback-failed-lab.json"

log "transactions prune dry-run report: $dry_report"
log "transactions prune report: $prune_report"
