#!/usr/bin/env bash
set -euo pipefail

# AI エージェントが使うための安全な git コミット & プッシュスクリプト
# 人間があらかじめリポジトリを clone 済み & リモート設定済みであることが前提。

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

REMOTE="${AUTOPILOT_GIT_REMOTE:-origin}"
BRANCH="${AUTOPILOT_GIT_BRANCH:-main}"

COMMIT_MSG="${1:-ai: autopilot update}"

# 変更がなければ何もしない
if git diff --quiet && git diff --cached --quiet; then
  echo "[ai-git] no changes to commit."
  exit 0
fi

echo "[ai-git] adding changes..."
git add -A

echo "[ai-git] committing..."
git commit -m "$COMMIT_MSG" || {
  echo "[ai-git] commit failed (maybe no changes?)."
  exit 0
}

echo "[ai-git] pushing to ${REMOTE} ${BRANCH}..."
git push "$REMOTE" "$BRANCH"
echo "[ai-git] done."
