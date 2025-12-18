# AI Autopilot State

このファイルは **AI エージェント自身が更新する進捗メモ** です。  
人間が初期化しても構いませんが、以降は codex / claude がここを書き換えながら開発を進めます。

---

## ゴール

Go 言語でシンプルな REST API サーバーを実装し、
`go test ./...` がすべて成功する状態まで持っていくこと。
必要な README やサンプルリクエストも整備すること。

---

## 現在の状況

- `projects/sample-go-api` を作成し、`/health` エンドポイントを提供するシンプルなサーバーを用意。
- `/health` を検証する最小限のテスト (`go test ./...`) が追加され、現状は緑。
- プロジェクト追加手順と命名ルールを `projects/README.md` に追記。

---

## TODO

- [x] プロジェクトのディレクトリ構成を確認する
- [x] エントリポイントとなる `main.go` を作成する
- [x] 最初のテストケースを追加する
- [x] テストが通るように実装を整える
- [ ] `/health` 以外のサンプルエンドポイントやエラーハンドリングを追加する（必要に応じて）
- [ ] README に今後の拡張方針や API 例を充実させる（必要に応じて）

---

## ログ

ここには AI エージェントが自由にメモしてよい領域です。

- 2025-12-18 08:52: `projects/sample-go-api` を作成し、`/health` エンドポイントとテストを追加。`projects/README.md` にプロジェクト追加手順と命名ルールを追記。`autopilot.yml` の `project_dir` は `./projects/sample-go-api` のままで問題ないことを確認。
- 2025-12-18 08:55: `bash ./scripts/ai-git-commit-push.sh \"ai: add sample go api project\"` を実行しコミット作成済み。リモートが未設定 (`git remote -v` で何も表示されず) のため push に失敗した。
