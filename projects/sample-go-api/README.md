# Sample Go API

シンプルなタスク管理用 REST API のサンプル実装です。標準ライブラリのみで構築されています。

## 動作要件
- Go 1.25 以上

## セットアップ
```
go mod tidy
```

## 実行方法
```
go run ./cmd/server
```

デフォルトでは `:8080` で起動します。ポートを変更する場合は `PORT` 環境変数を設定してください。

### サンプルリクエスト
- ヘルスチェック
  ```
  curl http://localhost:8080/health
  ```
- タスク作成
  ```
  curl -X POST http://localhost:8080/tasks \
    -H 'Content-Type: application/json' \
    -d '{"title":"Write tests"}'
  ```
- タスク一覧取得
  ```
  curl http://localhost:8080/tasks
  ```
- タスク更新
  ```
  curl -X PUT http://localhost:8080/tasks/1 \
    -H 'Content-Type: application/json' \
    -d '{"completed":true}'
  ```
- タスク削除
  ```
  curl -X DELETE http://localhost:8080/tasks/1
  ```

## テスト
```
go test ./...
```
