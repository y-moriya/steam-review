package i18n

func getEnglishMessages() map[string]string {
	return map[string]string{
		// Application info
		"app.name":    "Steam Reviews CLI Tool",
		"app.version": "Steam Reviews CLI Tool version %s",
		"app.started": "%s started",

		// Usage and help
		"usage.title":    "Usage:\n  steam-review [options]",
		"usage.options":  "Options:",
		"usage.examples": "Examples:",

		// Error messages
		"error.no_input":           "Error: Please specify either App ID or game name",
		"error.both_inputs":        "Error: Cannot specify both App ID and game name",
		"error.dir_creation":       "Failed to create output directory: %v",
		"error.review_fetch":       "Review fetch error: %v",
		"error.file_save":          "File save error: %v",
		"error.logger_init":        "Failed to initialize logger: %v",
		"error.game_details_fetch": "Failed to fetch game details: %v",

		// Success messages
		"success.completed":  "Process completed",
		"success.file_saved": "Reviews saved to %s",

		// Statistics
		"stats.title":              "=== Review Statistics ===",
		"stats.game":               "Game: %s",
		"stats.total_reviews":      "Total reviews: %d",
		"stats.positive":           "Positive: %d (%.1f%%)",
		"stats.negative":           "Negative: %d (%.1f%%)",
		"stats.language_breakdown": "Review Statistics by Language:",
		"stats.no_reviews":         "No reviews found",

		// File output
		"file.saved_files":         "=== Saved Files ===",
		"file.language_stats":      "  %s: %d reviews (%.1f%%) - Positive: %d (%.1f%%), Negative: %d",
		"file.creation_error":      "File creation error: %w",
		"file.game_details":        "=== Game Details ===",
		"file.game_name":           "Game Name: %s",
		"file.app_id":              "App ID: %s",
		"file.developer":           "Developer: %s",
		"file.publisher":           "Publisher: %s",
		"file.release_date":        "Release Date: %s",
		"file.price":               "Price: %s",
		"file.genres":              "Genres: %s",
		"file.categories":          "Categories: %s",
		"file.website":             "Website: %s",
		"file.age_restriction":     "Age Restriction: %d years and older",
		"file.free":                "Free: %t",
		"file.retrieved_at":        "Retrieved At: %s",
		"file.reviews_list":        "=== Reviews List ===",
		"file.review_number":       "=== Review %d ===",
		"file.json_write_error":    "JSON write error: %w",
		"file.language_save_error": "Language %s file save error: %v",
		"file.language_saved":      "Language %s: %d reviews saved to %s",
		"file.all_languages_saved": "All languages summary file saved: %s (%d reviews)",
		"file.summary_error":       "Summary file save error: %w",

		// Verbose logging
		"verbose.review_saved":   "Reviews saved to %s",
		"verbose.language_saved": "Language %s: %d reviews saved to %s",

		// Data fields (for output files)
		"field.developer":    "Developer",
		"field.publisher":    "Publisher",
		"field.release_date": "Release Date",
		"field.price":        "Price",
		"field.genre":        "Genre",
		"field.category":     "Category",
		"field.playtime":     "%d minutes",
		"field.review":       "review",
	}
}
