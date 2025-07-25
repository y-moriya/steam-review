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
	MsgUsageFull     = "usage.full_text"

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

	// API関連エラーメッセージ
	MsgErrorSteamAPIFetch    = "error.steam_api_fetch"
	MsgErrorJSONDecode       = "error.json_decode"
	MsgErrorGameNotFound     = "error.game_not_found"
	MsgErrorHTTPRequest      = "error.http_request"
	MsgErrorHTTPStatus       = "error.http_status"
	MsgErrorSteamAPIResponse = "error.steam_api_response"
	MsgErrorAppIDFetch       = "error.app_id_fetch"
	MsgErrorSteamStoreFetch  = "error.steam_store_fetch"
	MsgErrorAppDataNotFound  = "error.app_data_not_found"
	MsgErrorGameDetailsFail  = "error.game_details_fail"

	// Verbose/Progress ログメッセージ
	MsgVerboseReviewFetchStart    = "verbose.review_fetch_start"
	MsgVerboseReviewProgress      = "verbose.review_progress"
	MsgVerboseNoMoreReviews       = "verbose.no_more_reviews"
	MsgVerboseMaxReviewsReached   = "verbose.max_reviews_reached"
	MsgVerboseCursorNotChanged    = "verbose.cursor_not_changed"
	MsgVerboseTotalReviewsFetched = "verbose.total_reviews_fetched"
	MsgVerboseGameReviewFetch     = "verbose.game_review_fetch"
	MsgVerboseGameDetailsFetch    = "verbose.game_details_fetch"
	MsgVerboseGameDetailsObtained = "verbose.game_details_obtained"

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
