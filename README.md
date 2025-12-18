# AI Autopilot Template (codex / claude CLI 用)

このリポジトリは、

- **codex CLI**
- **claude CLI**

のような「AIエージェント付き CLI ツール」を前提とした  
**自律開発用テンプレート**です。

## コンセプト

- エンジニアは最初にゴールとルールを定義するだけ
- `scripts/codex-autopilot.sh` または `scripts/claude-autopilot.sh` を一度実行
- あとは AI エージェントが
  - プロジェクトのコード編集
  - テスト実行
  - リファクタリング
  - 進捗の記録（`state/autopilot_state.md`）
  を**自律的に繰り返す**
- クレジット切れやレートリミットで CLI が落ちた場合も、
  - `state/autopilot_state.md` に残った状態を読み直して
  - スクリプト側が再起動して続きを実行する

## 前提

- codex CLI / claude CLI がローカルで使えること（詳細な前提条件と設定例は後述）
- API キーや課金の設定は、それぞれの CLI ツール側で行うこと  
  （このリポジトリでは API キーは扱いません）

## ディレクトリ構成

```txt
.
├─ README.md
├─ autopilot.yml              # ゴールや対象プロジェクトの設定
├─ .codex/
│   └─ AGENTS.md              # codex 用エージェントプロンプト
├─ .claude/
│   └─ CLAUDE.md              # claude 用エージェントプロンプト
├─ state/
│   ├─ autopilot_state.md     # AI エージェントが自分で更新する進捗メモ
│   └─ last_session.log       # スクリプトが吐くログ
├─ scripts/
│   ├─ codex-autopilot.sh     # codex 用ウォッチャースクリプト
│   └─ claude-autopilot.sh    # claude 用ウォッチャースクリプト
└─ projects/
    └─ README.md              # プロジェクト配置に関する説明

## codex / claude CLI の前提

- codex CLI / claude CLI がローカルで使えること
- API キーや課金の設定は、それぞれの CLI ツール側で事前に済ませること  
  - 例: 環境変数（`CODEX_API_KEY` / `ANTHROPIC_API_KEY` など）や CLI 独自の設定コマンドで認証する
- `scripts/` 配下のシェルは実行権限を付与したうえで利用すること

## ウォッチャースクリプトの実行例

### codex-autopilot.sh

1. codex CLI をインストールし、API キー設定（例: `export CODEX_API_KEY=...`）を済ませる
2. プロジェクトルートで以下を実行

```bash
./scripts/codex-autopilot.sh
```

### claude-autopilot.sh

1. claude CLI をインストールし、API キー設定（例: `export ANTHROPIC_API_KEY=...`）を済ませる
2. プロジェクトルートで以下を実行

```bash
./scripts/claude-autopilot.sh
```

## サンプルプロジェクトの起動 / テスト

- `autopilot.yml` では `projects/sample-go-api` を対象とした Go REST API をゴールに設定しています
- プロジェクトを `projects/sample-go-api` に用意したら、以下のように実行・検証できます

```bash
cd projects/sample-go-api
go run ./...
go test ./...
```
