package stats

import (
	"fmt"

	"github.com/y-moriya/steam-review/internal/models"
)

// PrintReviewStats レビュー統計を表示
func PrintReviewStats(reviews []models.ReviewData, gameName string) {
	if len(reviews) == 0 {
		fmt.Printf("レビューが見つかりませんでした\n")
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

	fmt.Printf("\n=== レビュー統計 ===\n")
	fmt.Printf("ゲーム: %s\n", gameName)
	fmt.Printf("総レビュー数: %d\n", totalReviews)
	fmt.Printf("肯定的: %d (%.1f%%)\n", positiveReviews, positivePercent)
	fmt.Printf("否定的: %d (%.1f%%)\n", negativeReviews, negativePercent)

	fmt.Printf("\n言語別レビュー統計:\n")
	for lang, count := range languageCounts {
		positive := languagePositive[lang]
		negative := count - positive
		percent := float64(count) / float64(totalReviews) * 100
		positiveRate := float64(positive) / float64(count) * 100
		fmt.Printf("  %s: %d件 (%.1f%%) - 肯定的: %d件 (%.1f%%), 否定的: %d件\n",
			lang, count, percent, positive, positiveRate, negative)
	}
}
