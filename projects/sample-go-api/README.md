# sample-go-api

シンプルな Go 製 REST API サンプルです。デフォルトではポート `8080` で以下のエンドポイントを公開します。

- `GET /health` : 稼働確認用。`{"status":"ok"}` を返します。
- `POST /echo` : `{"message":"..."}` を受け取り、同じ内容を `echo` フィールドで返します。

## 実行方法

```bash
cd projects/sample-go-api
go run .
```

## テスト

```bash
cd projects/sample-go-api
go test ./...
```

## ディレクトリ構成

```txt
projects/sample-go-api
├─ go.mod
├─ main.go
└─ main_test.go
```
