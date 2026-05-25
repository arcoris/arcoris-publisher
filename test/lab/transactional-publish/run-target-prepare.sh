#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
require_cmd git
mkdir -p "$(lab_reports)" "$(lab_logs)" "$(lab_targets)"

report="$(lab_reports)/target-prepare.json"
stderr="$(lab_logs)/target-prepare.stderr"
preflight_report="$(lab_reports)/target-prepare-preflight.json"
preflight_stderr="$(lab_logs)/target-prepare-preflight.stderr"

run_step "seeding local bare remotes with main branches"
seed_bare_main arcoris/foundation
seed_bare_main arcoris/control

run_step "removing target worktrees so target prepare must clone them"
rm -rf "$(repository_worktree arcoris/foundation)" "$(repository_worktree arcoris/control)"

remote_template="$(target_prepare_remote_template)"
run_step "running target prepare with remote template: $remote_template"
code="$(capture_arcpub "$report" "$stderr" \
  target prepare \
  --manifest "$(lab_source)/arcpub.yaml" \
  --version v0.1.0 \
  --target-root "$(lab_targets)" \
  --remote-template "$remote_template" \
  --output json)"
assert_code "$code" 0 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"

for repo in arcoris/foundation arcoris/control; do
  worktree="$(repository_worktree "$repo")"
  bare="$(bare_repo "$repo")"
  [[ -d "$worktree/.git" ]] || fail "target worktree was not cloned: $worktree"
  [[ "$(git -C "$worktree" branch --show-current)" == "main" ]] || fail "target branch is not main: $worktree"
  [[ "$(git -C "$worktree" remote get-url origin)" == "file://$bare" ]] || fail "unexpected origin for $repo"
  assert_clean_worktree "$worktree"
done

if [[ -e "$(lab_state_dir)/publish.lock" || -d "$(lab_state_dir)/transactions" ]]; then
  fail "target prepare created publish transaction state"
fi

run_step "running preflight after target prepare"
code="$(run_common_arcpub_json "$preflight_report" "$preflight_stderr" preflight)"
assert_code "$code" 0 "$preflight_stderr"
assert_json_if_python "$preflight_report"
assert_no_path_leak "$preflight_report"

log "target prepare report: $report"
log "preflight-after-prepare report: $preflight_report"
