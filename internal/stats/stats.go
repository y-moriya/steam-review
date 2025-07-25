package stats

import (
	"github.com/y-moriya/steam-review/internal/models"
	"github.com/y-moriya/steam-review/pkg/i18n"
)

// Logger インターフェース（ロガーの依存性注入のため）
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// PrintReviewStats レビュー統計を表示
func PrintReviewStats(reviews []models.ReviewData, gameName string, logger Logger) {
	if len(reviews) == 0 {
		logger.Println(i18n.T(i18n.MsgStatsNoReviews))
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
	logger.Println(i18n.T(i18n.MsgStatsTitle))
	logger.Printf(i18n.Tf(i18n.MsgStatsGame, gameName))
	logger.Printf(i18n.Tf(i18n.MsgStatsTotalReviews, totalReviews))
	logger.Printf(i18n.Tf(i18n.MsgStatsPositive, positiveReviews, positivePercent))
	logger.Printf(i18n.Tf(i18n.MsgStatsNegative, negativeReviews, negativePercent))

	logger.Println()
	logger.Println(i18n.T(i18n.MsgStatsLanguageBreakdown))
	for lang, count := range languageCounts {
		positive := languagePositive[lang]
		negative := count - positive
		percent := float64(count) / float64(totalReviews) * 100
		positiveRate := float64(positive) / float64(count) * 100
		logger.Printf(i18n.Tf(i18n.MsgFileLanguageStats,
			lang, count, percent, positive, positiveRate, negative))
	}
}
