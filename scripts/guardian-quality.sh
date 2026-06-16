#!/usr/bin/env bash

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"

staged_files=()
while IFS= read -r file; do
  staged_files+=("$file")
done < <(git diff --cached --name-only --diff-filter=ACMR)

if [ "${#staged_files[@]}" -eq 0 ]; then
  echo "[Backend quality] No staged backend files."
  exit 0
fi

check_file_lengths() {
  local failed=0
  local file line_count

  for file in "${staged_files[@]}"; do
    case "$file" in
      *.md|*.sum|*.mod|*.json|*.lock|*.gitkeep)
        continue
        ;;
    esac

    if [ ! -f "$ROOT/$file" ]; then
      continue
    fi

    line_count="$(wc -l < "$ROOT/$file" | tr -d ' ')"
    if [ "$line_count" -gt 500 ]; then
      echo "[Backend quality] File too large: $file ($line_count lines, max 500)"
      failed=1
    fi
  done

  return "$failed"
}

check_function_lengths() {
  local failed=0
  local file

  for file in "${staged_files[@]}"; do
    case "$file" in
      *.go)
        ;;
      *)
        continue
        ;;
    esac

    if [ ! -f "$ROOT/$file" ]; then
      continue
    fi

    if ! awk -v limit=50 -v file="$file" '
      /^[[:space:]]*func[[:space:]]+/ {
        in_func = 1
        func_start = NR
        line = $0
        opens = gsub(/\{/, "{", line)
        closes = gsub(/\}/, "}", line)
        depth = opens - closes
        if (depth <= 0) {
          depth = 1
        }
        next
      }

      {
        if (in_func) {
          line = $0
          opens = gsub(/\{/, "{", line)
          closes = gsub(/\}/, "}", line)
          depth += opens - closes
          if (depth <= 0) {
            func_len = NR - func_start + 1
            if (func_len > limit) {
              printf("[Backend quality] Function too large: %s:%d (%d lines, max %d)\n", file, func_start, func_len, limit)
              failed = 1
            }
            in_func = 0
            depth = 0
          }
        }
      }

      END {
        exit failed
      }
    ' "$ROOT/$file"; then
      failed=1
    fi
  done

  return "$failed"
}

check_secrets() {
  local failed=0
  local file
  local secret_pattern

  secret_pattern='((OPENAI_API_KEY|SUPABASE_SERVICE_ROLE_KEY|DEEPGRAM_API_KEY|REDIS_URL)[[:space:]]*[:=][[:space:]]*["'"'"'A-Za-z0-9_\/+=.-]{8,}|-----BEGIN (RSA|EC|OPENSSH) PRIVATE KEY-----)'

  for file in "${staged_files[@]}"; do
    if [ ! -f "$ROOT/$file" ]; then
      continue
    fi

    if grep -En "$secret_pattern" "$ROOT/$file" >/dev/null; then
      echo "[Backend quality] Possible secret detected in $file"
      failed=1
    fi
  done

  return "$failed"
}

run_optional_tooling() {
  if [ ! -f "$ROOT/go.mod" ]; then
    echo "[Backend quality] go.mod not found. Skipping Go tooling checks."
    return 0
  fi

  if command -v go >/dev/null 2>&1; then
    gofmt_output="$(gofmt -l "$ROOT" 2>/dev/null || true)"
    if [ -n "$gofmt_output" ]; then
      echo "[Backend quality] gofmt required:"
      echo "$gofmt_output"
      return 1
    fi

    (cd "$ROOT" && go vet ./...)

    if command -v staticcheck >/dev/null 2>&1; then
      (cd "$ROOT" && staticcheck ./...)
    fi

    if command -v golangci-lint >/dev/null 2>&1; then
      (cd "$ROOT" && golangci-lint run)
    fi
  else
    echo "[Backend quality] Go toolchain not found. Skipping Go tooling checks."
  fi
}

check_file_lengths
check_function_lengths
check_secrets
run_optional_tooling

echo "[Backend quality] Passed."
