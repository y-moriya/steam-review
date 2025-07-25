package i18n

// メッセージキー定数
const (
	// アプリケーション情報
	MsgAppName    = "app.name"
	MsgAppVersion = "app.version"
	MsgAppStarted = "app.started"

	// 使用方法とヘルプ
	MsgUsageTitle    = "usage.title"
	MsgUsageOptions  = "usage.options"
	MsgUsageExamples = "usage.examples"
	MsgUsageHelp     = "usage.help_text"

	// エラーメッセージ
	MsgErrorNoInput         = "error.no_input"
	MsgErrorBothInputs      = "error.both_inputs"
	MsgErrorDirCreation     = "error.dir_creation"
	MsgErrorReviewFetch     = "error.review_fetch"
	MsgErrorFileSave        = "error.file_save"
	MsgErrorLoggerInit      = "error.logger_init"
	MsgErrorGameDetailsInit = "error.game_details_fetch"

	// 成功メッセージ
	MsgSuccessCompleted = "success.completed"
	MsgSuccessFileSaved = "success.file_saved"

	// 統計情報
	MsgStatsTitle             = "stats.title"
	MsgStatsGame              = "stats.game"
	MsgStatsTotalReviews      = "stats.total_reviews"
	MsgStatsPositive          = "stats.positive"
	MsgStatsNegative          = "stats.negative"
	MsgStatsLanguageBreakdown = "stats.language_breakdown"
	MsgStatsNoReviews         = "stats.no_reviews"

	// ファイル出力
	MsgFileSavedFiles        = "file.saved_files"
	MsgFileLanguageStats     = "file.language_stats"
	MsgFileCreationError     = "file.creation_error"
	MsgFileGameDetails       = "file.game_details"
	MsgFileGameName          = "file.game_name"
	MsgFileAppID             = "file.app_id"
	MsgFileDeveloper         = "file.developer"
	MsgFilePublisher         = "file.publisher"
	MsgFileReleaseDate       = "file.release_date"
	MsgFilePrice             = "file.price"
	MsgFileGenres            = "file.genres"
	MsgFileCategories        = "file.categories"
	MsgFileWebsite           = "file.website"
	MsgFileAgeRestriction    = "file.age_restriction"
	MsgFileFree              = "file.free"
	MsgFileRetrievedAt       = "file.retrieved_at"
	MsgFileReviewsList       = "file.reviews_list"
	MsgFileReviewNumber      = "file.review_number"
	MsgFileJSONWriteError    = "file.json_write_error"
	MsgFileLanguageSaveError = "file.language_save_error"
	MsgFileLanguageSaved     = "file.language_saved"
	MsgFileAllLanguagesSaved = "file.all_languages_saved"
	MsgFileSummaryError      = "file.summary_error"

	// 詳細ログ
	MsgVerboseReviewSaved   = "verbose.review_saved"
	MsgVerboseLanguageSaved = "verbose.language_saved"

	// データフィールド（出力ファイル用）
	MsgFieldDeveloper   = "field.developer"
	MsgFieldPublisher   = "field.publisher"
	MsgFieldReleaseDate = "field.release_date"
	MsgFieldPrice       = "field.price"
	MsgFieldGenre       = "field.genre"
	MsgFieldCategory    = "field.category"
	MsgFieldPlaytime    = "field.playtime"
	MsgFieldReview      = "field.review"
)
