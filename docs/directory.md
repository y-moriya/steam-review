# Steam Reviews CLI - 推奨ファイル構造

## ディレクトリ構造

```
steam-reviews-cli/
├── cmd/
│   └── steam-review/
│       └── main.go              # エントリーポイント、CLI引数処理
├── internal/
│   ├── api/
│   │   └── steam.go             # Steam API関連の処理
│   ├── models/
│   │   └── review.go            # データ構造体定義
│   ├── storage/
│   │   └── file.go              # ファイル保存処理
│   └── stats/
│       └── stats.go             # 統計処理
├── pkg/
│   └── config/
│       └── config.go            # 設定関連（外部から利用可能）
├── output/                      # デフォルトのファイル出力先
├── test_output/                 # テスト用のファイル出力先
├── go.mod
├── go.sum
└── README.md
```

## 各ファイルの役割

### `cmd/steam-review/main.go`
- アプリケーションのエントリーポイント
- コマンドライン引数の解析
- 各パッケージの協調処理
- エラーハンドリングとログ出力

### `internal/models/review.go`
- データ構造体の定義
- JSON変換ロジック
- データ変換メソッド

### `internal/api/steam.go`
- Steam API との通信
- レビューデータの取得
- API レスポンスの処理
- レート制限対応

### `internal/storage/file.go`
- ファイル保存処理
- JSON/テキスト形式での出力
- 言語別ファイル分割
- ディレクトリ作成

### `internal/stats/stats.go`
- レビュー統計の計算
- 統計情報の表示
- 言語別分析

### `pkg/config/config.go`
- 設定構造体
- デフォルト値定義
- バリデーション

## ファイル分割の利点

### 1. **保守性の向上**
- 機能ごとに分離されているため、修正箇所が特定しやすい
- 単一責任の原則に従い、各ファイルが明確な役割を持つ

### 2. **テスタビリティ**
- 各パッケージを独立してテストできる
- モックやスタブを使った単体テストが書きやすい

### 3. **再利用性**
- `pkg/` 以下は他のプロジェクトからも利用可能
- `internal/` は内部実装として隠蔽

### 4. **並行開発**
- 複数の開発者が異なるファイルを同時に編集可能
- Git でのマージコンフリクトが起きにくい

## パッケージ設計の指針

### `internal/` vs `pkg/`
- `internal/`: アプリケーション固有のロジック、外部からアクセス不可
- `pkg/`: 他のプロジェクトでも再利用可能なライブラリコード

### 依存関係の方向
```
cmd/ → internal/ → pkg/
```
- `cmd` は `internal` と `pkg` を利用
- `internal` は `pkg` を利用
- 循環依存は避ける

### インターフェース設計
各パッケージ間はインターフェースで結合し、具体的な実装に依存しない設計にする

## 追加の改善案

### 1. テストファイルの追加
```
internal/api/steam_test.go
internal/storage/file_test.go
internal/stats/stats_test.go
```

### 2. 設定ファイル対応
```
configs/
└── default.yaml
```

### 3. ドキュメント整備
```
docs/
├── API.md
├── USAGE.md
└── DEVELOPMENT.md
```

### 4. CI/CD設定
```
.github/
└── workflows/
    ├── test.yml
    └── release.yml
```

この構造により、コードの可読性、保守性、テスタビリティが大幅に向上し、将来的な機能拡張も容易になります。