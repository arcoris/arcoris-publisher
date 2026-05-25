#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

require_cmd git
require_cmd go

run_step "resetting lab root: $(lab_root)"
reset_lab_root

run_step "copying local publish fixture"
write_fixture_source
init_git_repo "$(lab_source)"

run_step "creating local bare remotes"
init_bare_repo "$(bare_repo arcoris/foundation)"
init_bare_repo "$(bare_repo arcoris/control)"

run_step "creating target worktrees"
init_target_worktree "$(repository_worktree arcoris/foundation)" "$(bare_repo arcoris/foundation)"
init_target_worktree "$(repository_worktree arcoris/control)" "$(bare_repo arcoris/control)"

run_step "building arcpub"
bash "$LAB_DIR/build.sh"

log "lab ready: $(lab_root)"
log "next: bash $LAB_DIR/run-happy-path.sh"
