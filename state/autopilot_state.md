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

- `projects/sample-go-api` を作成し、シンプルな REST API（/health, /echo）とテストを追加。
- `projects/README.md` と `README.md` にプロジェクト配置・命名規則、`project_dir` 設定方法、構成例を追記。

---

## TODO

- [ ] API の追加要件があれば拡張する（例: ルーティングやエラーハンドリングの充実）
- [ ] README に API 例やクライアントサンプルを増やす
- [ ] CI など自動テスト環境が必要なら整備する

---

## ログ

- 2025-02-03: `projects/sample-go-api` を作成し、/health と /echo エンドポイントのテストを追加。プロジェクト配置と命名規則のドキュメントを整備。
