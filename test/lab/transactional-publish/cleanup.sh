#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

guard_lab_root
run_step "removing lab root: $(lab_root)"
rm -rf "$(lab_root)"
log "removed"
