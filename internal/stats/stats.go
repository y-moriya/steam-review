package stats

import (
	"github.com/y-moriya/steam-review/internal/models"
)

// Logger インターフェース（ロガーの依存性注入のため）
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// PrintReviewStats レビュー統計を表示
func PrintReviewStats(reviews []models.ReviewData, gameName string, logger Logger) {
	if len(reviews) == 0 {
		logger.Println("レビューが見つかりませんでした")
		return
	}

	totalReviews := len(reviews)
	positiveReviews := 0
	languageCounts := make(map[string]int)
	languagePositive := make(map[string]int)

	for _, review := range reviews {
		if review.VotedUp {
			positiveReviews++
		}
		lang := review.Language
		if lang == "" {
			lang = "unknown"
		}
		languageCounts[lang]++
		if review.VotedUp {
			languagePositive[lang]++
		}
	}

	negativeReviews := totalReviews - positiveReviews
	positivePercent := float64(positiveReviews) / float64(totalReviews) * 100
	negativePercent := float64(negativeReviews) / float64(totalReviews) * 100

	logger.Println()
	logger.Println("=== レビュー統計 ===")
	logger.Printf("ゲーム: %s", gameName)
	logger.Printf("総レビュー数: %d", totalReviews)
	logger.Printf("肯定的: %d (%.1f%%)", positiveReviews, positivePercent)
	logger.Printf("否定的: %d (%.1f%%)", negativeReviews, negativePercent)

	logger.Println()
	logger.Println("言語別レビュー統計:")
	for lang, count := range languageCounts {
		positive := languagePositive[lang]
		negative := count - positive
		percent := float64(count) / float64(totalReviews) * 100
		positiveRate := float64(positive) / float64(count) * 100
		logger.Printf("  %s: %d件 (%.1f%%) - 肯定的: %d件 (%.1f%%), 否定的: %d件",
			lang, count, percent, positive, positiveRate, negative)
	}
}
