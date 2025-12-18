# Sample Go API

シンプルな Go 製 REST API サーバーのサンプルです。`/health` エンドポイントでアプリケーションの稼働状態を返します。

## セットアップ

1. Go 1.25 以降をインストールします。
2. 依存関係を取得します。
   ```bash
   go mod tidy
   ```

## 実行方法

ローカルでサーバーを起動します。

```bash
# プロジェクトルート (projects/sample-go-api) で実行
GO_RUN_ADDRESS=:8080 # 任意。環境変数は未使用だがポート指定の例です。
go run .
```

起動後、別のターミナルで疎通確認できます。

```bash
curl http://localhost:8080/health
```

期待されるレスポンス:

```json
{"status":"ok"}
```

## テスト

すべてのテストを実行します。

```bash
cd projects/sample-go-api
go test ./...
```

## 構成

- `main.go`: エントリーポイントとルーティングの定義。
- `health_test.go`: `/health` エンドポイントの動作確認テスト。
