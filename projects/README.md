# projects ディレクトリの使い方

このリポジトリで管理する各プロジェクトは、すべて `projects/` 配下に配置します。サンプルとして `projects/sample-go-api` を追加しています。

## 配置場所と命名規則

- プロジェクトは `projects/<project-name>/` に置きます。
- `<project-name>` は **kebab-case**（例: `sample-go-api`, `next-dashboard`）で統一してください。
- 各プロジェクトには最低限 `README.md` とテストを含め、`go`, `node`, `python` など使用言語に合わせた標準的な構成を守ります。
- `autopilot.yml` の `project_dir` には、対象プロジェクトへの **相対パス**（例: `./projects/sample-go-api`）を設定してください。

## 追加済みサンプルプロジェクト

- **sample-go-api**: Go 製のシンプルな REST API。`GET /health` と `POST /echo` を提供します。
- テスト実行: `cd projects/sample-go-api && go test ./...`

## 今後プロジェクトを増やすときの構成例

```txt
projects/
├─ sample-go-api/          # Go の REST API サンプル
│  ├─ README.md
│  ├─ go.mod
│  └─ main.go
├─ next-frontend/          # Next.js や React のフロントエンド
│  ├─ README.md
│  ├─ package.json
│  └─ src/
└─ data-pipeline/          # Python 製バッチや ETL
   ├─ README.md
   ├─ pyproject.toml
   └─ src/
```

この例にならい、言語ごとの標準ツールチェーンを利用し、プロジェクト固有のセットアップ手順やテスト手順は各 `README.md` に明記してください。
