package i18n

func getJapaneseMessages() map[string]string {
	return map[string]string{
		// アプリケーション情報
		"app.name":    "Steam Reviews CLI Tool",
		"app.version": "Steam Reviews CLI Tool version %s",

		// 使用方法とヘルプ
		"usage.title":    "使用方法:\n  steam-review [オプション]",
		"usage.options":  "オプション:",
		"usage.examples": "使用例:",

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
		"file.saved_files":    "=== 保存したファイル一覧 ===",
		"file.language_stats": "  %s: %d件 (%.1f%%) - 肯定的: %d件 (%.1f%%), 否定的: %d件",

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
