#!/usr/bin/env bash

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"

staged_files=()
while IFS= read -r file; do
  staged_files+=("$file")
done < <(git diff --cached --name-only --diff-filter=ACMR)

if [ "${#staged_files[@]}" -eq 0 ]; then
  echo "[Backend behavior] No staged backend files."
  exit 0
fi

for required_dir in internal cmd tests; do
  if [ ! -d "$ROOT/$required_dir" ]; then
    echo "[Backend behavior] Missing required directory: $required_dir"
    exit 1
  fi
done

if [ -f "$ROOT/go.mod" ] && command -v go >/dev/null 2>&1; then
  (cd "$ROOT" && go test ./...)
else
  echo "[Backend behavior] go.mod or Go toolchain missing. Skipping runtime checks until backend is initialized."
fi

echo "[Backend behavior] Passed."
