package i18n

func getJapaneseMessages() map[string]string {
	return map[string]string{
		// アプリケーション情報
		"app.name":    "Steam Reviews CLI Tool",
		"app.version": "Steam Reviews CLI Tool version %s",
		"app.started": "%s が開始されました",

		// 使用方法とヘルプ
		"usage.title":    "使用方法:\n  steam-review [オプション]",
		"usage.options":  "オプション:",
		"usage.examples": "使用例:",
		"usage.full_text": `%s version %s

使用方法:
  steam-review [オプション]

オプション:
  -appid string         Steam App ID (例: 440)
  -game string          ゲーム名 (例: "Team Fortress 2")
  -max int             最大取得レビュー数 (デフォルト: 100, 0で無制限)
  -lang string         取得する言語 (カンマ区切り, デフォルト: japanese, 例: "japanese,english")
  -output string       出力ディレクトリ (デフォルト: output)
  -split              言語別にファイルを分けて保存
  -json               出力ファイルをJSON形式(.json)にする (デフォルト: テキスト形式)
  -verbose            詳細なログを表示
  -filter string      レビューのフィルター (recent: 作成日時順, updated: 更新日時順, all: 有用性順(デフォルト))
  -help               このヘルプを表示
  -version            バージョン情報を表示

使用例:
  # App IDを指定して日本語レビューを取得（デフォルト: 有用性順）
  steam-review -appid 440 -max 500 -verbose

  # 作成日時順でレビューを取得
  steam-review -appid 440 -max 500 -filter recent -verbose

  # ゲーム名で英語レビューを取得
  steam-review -game "Cyberpunk 2077" -lang "english" -max 1000 -output ./reviews

  # 複数言語のレビューを取得
  steam-review -game "Elden Ring" -lang "japanese,english" -max 300 -split

  # 日本語レビューをJSON形式で保存
  steam-review -appid 570 -max 2000 -output ./dota2_reviews -json -verbose

  # すべての言語のレビューを取得
  steam-review -appid 730 -lang "all" -max 1000 -split

  # 最近更新されたレビューから取得
  steam-review -appid 730 -filter updated -max 200

注意:
  - App IDとゲーム名のどちらか一方を指定してください
  - -lang を指定しない場合、デフォルトで日本語レビューのみを取得します
  - "all" を指定するとすべての言語のレビューを取得します
  - 大量のレビューを取得する場合は時間がかかります
  - Steam APIのレート制限により、リクエスト間に1秒の待機時間があります`,

		// エラーメッセージ
		"error.no_input":           "エラー: App ID またはゲーム名を指定してください",
		"error.both_inputs":        "エラー: App ID とゲーム名の両方を指定することはできません",
		"error.dir_creation":       "出力ディレクトリの作成に失敗しました: %v",
		"error.review_fetch":       "レビュー取得エラー: %v",
		"error.file_save":          "ファイル保存エラー: %v",
		"error.logger_init":        "ロガーの初期化に失敗しました: %v",
		"error.game_details_fetch": "ゲーム詳細情報の取得に失敗しました: %v",

		// 成功メッセージ
		"success.completed":  "処理が完了しました",
		"success.file_saved": "レビューを %s に保存しました",

		// 統計情報
		"stats.title":              "=== レビュー統計 ===",
		"stats.game":               "ゲーム: %s",
		"stats.total_reviews":      "総レビュー数: %d",
		"stats.positive":           "肯定的: %d (%.1f%%)",
		"stats.negative":           "否定的: %d (%.1f%%)",
		"stats.language_breakdown": "言語別レビュー統計:",
		"stats.no_reviews":         "レビューが見つかりませんでした",

		// ファイル出力
		"file.saved_files":         "=== 保存したファイル一覧 ===",
		"file.language_stats":      "  %s: %d件 (%.1f%%) - 肯定的: %d件 (%.1f%%), 否定的: %d件",
		"file.creation_error":      "ファイル作成エラー: %w",
		"file.game_details":        "=== ゲーム詳細情報 ===",
		"file.game_name":           "ゲーム名: %s",
		"file.app_id":              "App ID: %s",
		"file.developer":           "開発者: %s",
		"file.publisher":           "パブリッシャー: %s",
		"file.release_date":        "リリース日: %s",
		"file.price":               "価格: %s",
		"file.genres":              "ジャンル: %s",
		"file.categories":          "カテゴリ: %s",
		"file.website":             "ウェブサイト: %s",
		"file.age_restriction":     "年齢制限: %d歳以上",
		"file.free":                "無料: %t",
		"file.retrieved_at":        "情報取得日時: %s",
		"file.reviews_list":        "=== レビュー一覧 ===",
		"file.review_number":       "=== レビュー %d ===",
		"file.json_write_error":    "JSON書き込みエラー: %w",
		"file.language_save_error": "言語 %s のファイル保存エラー: %v",
		"file.language_saved":      "言語 %s: %d件のレビューを %s に保存",
		"file.all_languages_saved": "全言語統合ファイルを保存: %s (%d件)",
		"file.summary_error":       "サマリーファイル保存エラー: %w",

		// API関連エラーメッセージ
		"error.steam_api_fetch":    "Steam API取得エラー: %w",
		"error.json_decode":        "JSONデコードエラー: %w",
		"error.game_not_found":     "ゲーム '%s' が見つかりません",
		"error.http_request":       "HTTP リクエストエラー: %w",
		"error.http_status":        "HTTP エラー: %d",
		"error.steam_api_response": "Steam API エラー: success = %d",
		"error.app_id_fetch":       "App ID取得エラー: %w",
		"error.steam_store_fetch":  "Steam Store API取得エラー: %w",
		"error.app_data_not_found": "App ID %s のデータが見つかりません",
		"error.game_details_fail":  "App ID %s の詳細情報取得に失敗しました",

		// Verbose/Progress ログメッセージ
		"verbose.review_fetch_start":    "App ID %s のレビュー取得を開始します",
		"verbose.review_progress":       "現在のレビュー数: %d, カーソル: %s",
		"verbose.no_more_reviews":       "これ以上レビューがありません",
		"verbose.max_reviews_reached":   "最大レビュー数 %d に到達しました",
		"verbose.cursor_not_changed":    "カーソルが変更されませんでした。終了します",
		"verbose.total_reviews_fetched": "合計 %d 件のレビューを取得しました",
		"verbose.game_review_fetch":     "ゲーム '%s' (App ID: %s) のレビューを取得します",
		"verbose.game_details_fetch":    "App ID %s のゲーム詳細情報を取得中...",
		"verbose.game_details_obtained": "ゲーム詳細情報を取得しました: %s",

		// 詳細ログ
		"verbose.review_saved":   "レビューを %s に保存しました",
		"verbose.language_saved": "言語 %s: %d件のレビューを %s に保存",

		// データフィールド（出力ファイル用）
		"field.developer":    "開発者",
		"field.publisher":    "パブリッシャー",
		"field.release_date": "リリース日",
		"field.price":        "価格",
		"field.genre":        "ジャンル",
		"field.category":     "カテゴリ",
		"field.playtime":     "%d分",
		"field.review":       "review",
	}
}
