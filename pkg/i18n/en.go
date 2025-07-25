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
		"usage.full_text": `%s version %s

Usage:
  steam-review [options]

Options:
  -appid string         Steam App ID (e.g., 440)
  -game string          Game name (e.g., "Team Fortress 2")
  -max int             Maximum number of reviews to retrieve (default: 100, 0 for unlimited)
  -lang string         Languages to retrieve (comma-separated, default: japanese, e.g., "japanese,english")
  -output string       Output directory (default: output)
  -split              Split files by language
  -json               Output files in JSON format (.json) (default: text format)
  -verbose            Show detailed logs
  -filter string      Review filter (recent: by creation date, updated: by update date, all: by helpfulness (default))
  -help               Show this help
  -version            Show version information

Examples:
  # Get Japanese reviews by App ID (default: sorted by helpfulness)
  steam-review -appid 440 -max 500 -verbose

  # Get reviews sorted by creation date
  steam-review -appid 440 -max 500 -filter recent -verbose

  # Get English reviews by game name
  steam-review -game "Cyberpunk 2077" -lang "english" -max 1000 -output ./reviews

  # Get reviews in multiple languages
  steam-review -game "Elden Ring" -lang "japanese,english" -max 300 -split

  # Save Japanese reviews in JSON format
  steam-review -appid 570 -max 2000 -output ./dota2_reviews -json -verbose

  # Get reviews in all languages
  steam-review -appid 730 -lang "all" -max 1000 -split

  # Get recently updated reviews
  steam-review -appid 730 -filter updated -max 200

Notes:
  - Specify either App ID or game name, not both
  - If -lang is not specified, only Japanese reviews will be retrieved by default
  - Use "all" to retrieve reviews in all languages
  - Retrieving a large number of reviews may take time
  - Due to Steam API rate limits, there is a 1-second delay between requests`,

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
