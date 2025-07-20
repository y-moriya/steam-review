# Steam Reviews CLI Tool

Steam ReviewsはSteamゲームのレビューを取得・保存するためのコマンドラインツールです。
App IDまたはゲーム名を指定して、レビューを取得し、JSON形式またはテキスト形式で保存できます。

## 機能

- App IDまたはゲーム名でゲームを指定
- 言語でレビューをフィルタリング
- 最大取得件数の指定
- レビューの作成日時/更新日時でのソート
- JSON形式またはテキスト形式での保存
- 言語別のファイル分割
- 詳細な統計情報の表示

## インストール

```bash
go install github.com/y-moriya/steam-review@latest
```

## 使用方法

```
steam-review [オプション]
```

### オプション

| オプション | 説明 | デフォルト値 |
|------------|------|--------------|
| -appid     | Steam App ID (例: 440) | - |
| -game      | ゲーム名 (例: "Team Fortress 2") | - |
| -max       | 最大取得レビュー数 (0で無制限) | 100 |
| -lang      | 取得する言語 (カンマ区切り) | japanese |
| -output    | 出力ディレクトリ | output |
| -split     | 言語別にファイルを分けて保存 | false |
| -json      | 出力ファイルをJSON形式(.json)にする | false |
| -filter    | レビューのフィルター (recent/updated/all) | all |
| -verbose   | 詳細なログを表示 | false |
| -help      | ヘルプを表示 | false |
| -version   | バージョン情報を表示 | - |

### フィルターオプション

- `all`: 有用性による並び替え（デフォルト）
- `recent`: 作成日時による並び替え
- `updated`: 最終更新日時による並び替え

## 使用例

1. App IDを指定して日本語レビューを取得（デフォルト）
```bash
steam-review -appid 440 -max 500 -verbose
```

2. ゲーム名で英語レビューを取得
```bash
steam-review -game "Cyberpunk 2077" -lang "english" -max 1000 -output ./reviews
```

3. 複数言語のレビューを取得
```bash
steam-review -game "Elden Ring" -lang "japanese,english" -max 300 -split
```

4. 日本語レビューをJSON形式で保存
```bash
steam-review -appid 570 -max 2000 -output ./dota2_reviews -json -verbose
```

5. すべての言語のレビューを取得
```bash
steam-review -appid 730 -lang "all" -max 1000 -split
```

6. 最近更新されたレビューから取得
```bash
steam-review -appid 730 -filter updated -max 200
```

## 出力ファイル

### テキスト形式 (デフォルト)

```
=== レビュー 1 ===
ID: 195635539
language: japanese
voted_up: true
votes_up: 54
votes_funny: 28
weighted_score: 0.82
steam_purchase: true
playtime: 3754分
created_at: 2025-05-26 00:49:30
updated_at: 2025-05-26 22:12:04
review:
レビュー本文
```

### JSON形式 (-json オプション)

```json
[
  {
    "recommendation_id": "12345678",
    "author": {
      "steam_id": "76561197960287930",
      "num_games_owned": 100,
      "num_reviews": 10,
      "playtime_forever": 1000,
      "playtime_last_two_weeks": 10,
      "playtime_at_review": 500,
      "last_played": 1626825600
    },
    "language": "japanese",
    "review": "レビュー本文",
    "timestamp_created": 1626825600,
    "timestamp_updated": 1626825600,
    "voted_up": true,
    "votes_up": 10,
    "votes_funny": 5,
    "weighted_vote_score": 0.8,
    "comment_count": 3,
    "steam_purchase": true,
    "received_for_free": false,
    "written_during_early_access": false
  }
]
```

## 注意事項

- App IDとゲーム名のどちらか一方を指定してください
- `-lang` を指定しない場合、デフォルトで日本語レビューのみを取得します
- `all` を指定するとすべての言語のレビューを取得します
- 大量のレビューを取得する場合は時間がかかります
- Steam APIのレート制限により、リクエスト間に1秒の待機時間があります
- 出力ディレクトリは自動的に作成されます

## ライセンス

[MITライセンス](LICENSE)
