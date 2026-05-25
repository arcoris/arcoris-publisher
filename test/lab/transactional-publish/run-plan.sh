#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
mkdir -p "$(lab_reports)"
report="$(lab_reports)/plan.json"
stderr="$(lab_logs)/plan.stderr"

run_step "running plan"
code="$(capture_arcpub "$report" "$stderr" plan --manifest "$(lab_source)/arcpub.yaml" --version v0.1.0 --output json)"
assert_code "$code" 0 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"
log "report: $report"
