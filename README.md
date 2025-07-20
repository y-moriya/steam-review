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
| -txt       | 出力ファイルをテキスト形式(.txt)にする | false |
| -sort      | レビューのソート順 | created_desc |
| -verbose   | 詳細なログを表示 | false |
| -help      | ヘルプを表示 | false |

### ソートオプション

- `created_desc`: 作成日時の降順（新しい順）
- `created_asc`: 作成日時の昇順（古い順）
- `updated_desc`: 更新日時の降順（新しい更新順）
- `updated_asc`: 更新日時の昇順（古い更新順）

## 使用例

1. App IDを指定して日本語レビューを取得（デフォルト）
```bash
steam-review -appid 440 -max 500 -verbose
```

2. ゲーム名で英語レビューを取得
```bash
steam-review -game "Cyberpunk 2077" -lang "english" -max 1000 -output ./reviews
```

3. 複数言語のレビューをテキスト形式で取得
```bash
steam-review -game "Elden Ring" -lang "japanese,english" -max 300 -split -txt
```

4. 日本語レビューのみをテキスト形式で保存
```bash
steam-review -appid 570 -max 2000 -output ./dota2_reviews -txt -verbose
```

5. すべての言語のレビューを取得
```bash
steam-review -appid 730 -lang "all" -max 1000 -split
```

6. 最近更新されたレビューから取得
```bash
steam-review -appid 730 -sort updated_desc -max 200
```

## 出力ファイル

### JSON形式 (デフォルト)

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

### テキスト形式 (-txt オプション)

```
=== レビュー 1 ===
ID: 12345678
言語: japanese
評価: 肯定的
投票数: 10
面白い投票: 5
重み付けスコア: 0.80
Steam購入: true
プレイ時間: 500分
作成日時: 2023-07-21 00:00:00
更新日時: 2023-07-21 00:00:00
レビュー内容:
レビュー本文
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
