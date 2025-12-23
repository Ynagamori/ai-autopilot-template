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

- codex CLI / claude CLI がローカルで使えること
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
