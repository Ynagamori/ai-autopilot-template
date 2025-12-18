#!/usr/bin/env bash
set -euo pipefail

# claude CLI 用 自律開発ウォッチャースクリプト
# ※ 実際の claude CLI のオプションは各自の環境に合わせて修正してください。

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

  # TODO: あなたの環境の claude CLI に合わせて、以下のコマンドを変更してください。
  # 例:
  #   - プロジェクトルートとして $ROOT_DIR を指定
  #   - エージェント設定として .claude/AGENTS.md を読み込む
  #
  # 下の行はダミーです。実際の CLI 仕様に合わせて書き換えてください。
  claude \
    --project "$ROOT_DIR" \
    --agent-file "$ROOT_DIR/.claude/AGENTS.md" \
    2>&1 | tee -a "$LOG_FILE"

  EXIT_CODE=${PIPESTATUS[0]}
  echo "[autopilot] claude exited with code $EXIT_CODE" | tee -a "$LOG_FILE"

  if grep -qi "insufficient quota" "$LOG_FILE" || grep -qi "rate limit" "$LOG_FILE"; then
    echo "[autopilot] seems like credits exhausted or rate limited. sleeping 30min..." | tee -a "$LOG_FILE"
    sleep 1800
    continue
  fi

  exit "$EXIT_CODE"
done
