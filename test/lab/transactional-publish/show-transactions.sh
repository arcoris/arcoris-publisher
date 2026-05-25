#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
mkdir -p "$(lab_reports)"

run_step "listing transactions"
"$(lab_arcpub)" transactions list --state-dir "$(lab_state_dir)" --output text | tee "$(lab_reports)/transactions-list.txt"
"$(lab_arcpub)" transactions list --state-dir "$(lab_state_dir)" --output json >"$(lab_reports)/transactions-list.json"
assert_json_if_python "$(lab_reports)/transactions-list.json"
assert_no_path_leak "$(lab_reports)/transactions-list.json"

tx_id="${1:-$(latest_transaction_id || true)}"
if [[ -n "$tx_id" ]]; then
  run_step "showing transaction $tx_id"
  "$(lab_arcpub)" transactions show "$tx_id" --state-dir "$(lab_state_dir)" --output text | tee "$(lab_reports)/transaction-show.txt"
  "$(lab_arcpub)" transactions show "$tx_id" --state-dir "$(lab_state_dir)" --output json >"$(lab_reports)/transaction-show.json"
  assert_json_if_python "$(lab_reports)/transaction-show.json"
  assert_no_path_leak "$(lab_reports)/transaction-show.json"
else
  log "no transaction id found"
fi
