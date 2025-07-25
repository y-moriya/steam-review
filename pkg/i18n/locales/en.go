package i18n

func getEnglishMessages() map[string]string {
	return map[string]string{
		// Application info
		"app.name":    "Steam Reviews CLI Tool",
		"app.version": "Steam Reviews CLI Tool version %s",

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
		"file.saved_files":    "=== Saved Files ===",
		"file.language_stats": "  %s: %d reviews (%.1f%%) - Positive: %d (%.1f%%), Negative: %d",

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
