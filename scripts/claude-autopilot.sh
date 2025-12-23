#!/usr/bin/env bash
set -euo pipefail

# claude CLI 用 自律開発ウォッチャースクリプト
# Claude CLI の自動エージェント実行に合わせたオプションを指定する。

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

MAX_HOURS=24
START_TS=$(date +%s)

LOG_DIR="$ROOT_DIR/state"
mkdir -p "$LOG_DIR"

LOG_FILE="$LOG_DIR/last_session.log"

echo "[autopilot] claude session started at $(date)" | tee -a "$LOG_FILE"

while true; do
  NOW=$(date +%s)
  ELAPSED=$(( (NOW - START_TS) / 3600 ))

  if [ "$ELAPSED" -ge "$MAX_HOURS" ]; then
    echo "[autopilot] max hours reached (${ELAPSED}h), exiting." | tee -a "$LOG_FILE"
    exit 0
  fi

  echo "[autopilot] launching claude CLI..." | tee -a "$LOG_FILE"

  # 非対話モードの `claude run` を使い、
  # エージェント設定は .claude/CLAUDE.md から読み込む。
  set +e
  claude run \
    --project-dir "$ROOT_DIR" \
    --agent-file "$ROOT_DIR/.claude/CLAUDE.md" \
    2>&1 | tee -a "$LOG_FILE"
  EXIT_CODE=${PIPESTATUS[0]}
  set -e

  echo "[autopilot] claude exited with code $EXIT_CODE" | tee -a "$LOG_FILE"

  if grep -qi "insufficient quota" "$LOG_FILE" || grep -qi "rate limit" "$LOG_FILE"; then
    echo "[autopilot] seems like credits exhausted or rate limited. sleeping 30min..." | tee -a "$LOG_FILE"
    sleep 1800
    continue
  fi

  if [ "$EXIT_CODE" -ne 0 ]; then
    echo "[autopilot] non-zero exit code ($EXIT_CODE), stopping watcher." | tee -a "$LOG_FILE"
    exit "$EXIT_CODE"
  fi

  echo "[autopilot] claude finished one iteration; restarting for next task..." | tee -a "$LOG_FILE"
done
