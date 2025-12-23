#!/usr/bin/env bash
set -euo pipefail

# AI エージェントが使うための安全な git コミット & プッシュスクリプト
# 人間があらかじめリポジトリを clone 済み & リモート設定済みであることが前提。

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

REMOTE="${AUTOPILOT_GIT_REMOTE:-origin}"

TASK_NAME_RAW="${AUTOPILOT_TASK_NAME:-}"
if [ -n "$TASK_NAME_RAW" ]; then
  TASK_NAME=$(echo "$TASK_NAME_RAW" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9._-]/-/g; s/--*/-/g; s/^-//; s/-$//')
else
  TASK_NAME=""
fi

DEFAULT_BRANCH="feature/autopilot"
if [ -n "$TASK_NAME" ]; then
  DEFAULT_BRANCH="feature/${TASK_NAME}"
fi

BRANCH="${AUTOPILOT_GIT_BRANCH:-$DEFAULT_BRANCH}"

COMMIT_MSG="${1:-ai: autopilot update}"

echo "[ai-git] ensuring branch $BRANCH exists..."
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "$BRANCH" ]; then
  if git show-ref --verify --quiet "refs/heads/$BRANCH"; then
    git checkout "$BRANCH"
  else
    git checkout -b "$BRANCH"
  fi
fi

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
git push -u "$REMOTE" "$BRANCH"
echo "[ai-git] done."
