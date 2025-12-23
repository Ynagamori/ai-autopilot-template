# Sample Go Task API

シンプルな Go 製 REST API です。メモリ内にタスクを管理し、以下のエンドポイントを提供します。

- `GET /health` : ヘルスチェック
- `GET /tasks` : タスク一覧（`offset` と `limit` で件数制御可能）
- `POST /tasks` : タスク作成（例: `{"title": "Write tests"}`）
- `PATCH /tasks/{id}` : タスクタイトル更新（例: `{"title": "Ship features"}`）
- `PATCH /tasks/{id}/complete` : タスク完了
- `DELETE /tasks/{id}` : タスク削除

## 動作条件

- Go 1.22 以上

## セットアップ

プロジェクト直下で依存関係を解決します。

```bash
cd projects/sample-go-api
go mod tidy
```

## サーバー起動

```bash
PORT=8080 go run .
```

### サンプルリクエスト

```bash
# ヘルスチェック
curl http://localhost:8080/health

# タスク作成
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Write tests"}'

# 一覧取得
curl http://localhost:8080/tasks

# 一覧取得（2件目から1件だけ）
curl "http://localhost:8080/tasks?offset=1&limit=1"

# タイトル更新
curl -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Ship features"}'

# 完了マーク
curl -X PATCH http://localhost:8080/tasks/1/complete

# 削除
curl -X DELETE http://localhost:8080/tasks/1
```

## テスト

```bash
go test ./...
```
