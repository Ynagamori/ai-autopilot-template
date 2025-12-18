#!/usr/bin/env bash
set -euo pipefail

# codex CLI 用 自律開発ウォッチャースクリプト
# 依存: codex CLI がインストールされ、API キー設定（例: CODEX_API_KEY 環境変数や CLI 設定）が済んでいること
# 実行例:
#   chmod +x scripts/codex-autopilot.sh
#   ./scripts/codex-autopilot.sh
# ※ 実際の codex CLI のオプションは各自の環境に合わせて修正してください。

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

MAX_HOURS=24
START_TS=$(date +%s)

LOG_DIR="$ROOT_DIR/state"
mkdir -p "$LOG_DIR"

LOG_FILE="$LOG_DIR/last_session.log"

echo "[autopilot] codex session started at $(date)" | tee -a "$LOG_FILE"

while true; do
  NOW=$(date +%s)
  ELAPSED=$(( (NOW - START_TS) / 3600 ))

  if [ "$ELAPSED" -ge "$MAX_HOURS" ]; then
    echo "[autopilot] max hours reached (${ELAPSED}h), exiting." | tee -a "$LOG_FILE"
    exit 0
  fi

  echo "[autopilot] launching codex CLI..." | tee -a "$LOG_FILE"

  # 非対話モードの `codex exec` を使い、
  # AGENT 設定は .codex/AGENTS.md から読み込む。
  set +e
  codex exec \
    --full-auto \
    --cd "$ROOT_DIR" \
    - < "$ROOT_DIR/.codex/AGENTS.md" \
    2>&1 | tee -a "$LOG_FILE"
  EXIT_CODE=${PIPESTATUS[0]}
  set -e

  echo "[autopilot] codex exited with code $EXIT_CODE" | tee -a "$LOG_FILE"

  # ログにクレジット系エラーっぽい文言がないか簡易チェック
  if grep -qi "insufficient quota" "$LOG_FILE" || grep -qi "rate limit" "$LOG_FILE"; then
    echo "[autopilot] seems like credits exhausted or rate limited. sleeping 30min..." | tee -a "$LOG_FILE"
    sleep 1800
    continue
  fi

  # エラー終了コードの場合はそこで止める（人間の確認が必要な想定）
  if [ "$EXIT_CODE" -ne 0 ]; then
    echo "[autopilot] non-zero exit code ($EXIT_CODE), stopping watcher." | tee -a "$LOG_FILE"
    exit "$EXIT_CODE"
  fi

  # 正常終了コード(0)の場合は、max_hours に達するまで次のループを継続する
  echo "[autopilot] codex finished one iteration; restarting for next task..." | tee -a "$LOG_FILE"
done
