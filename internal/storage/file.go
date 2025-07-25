package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/y-moriya/steam-review/internal/models"
	"github.com/y-moriya/steam-review/pkg/config"
)

// SaveReviewsToFile レビューをファイルに保存
func SaveReviewsToFile(reviews []models.ReviewData, filename string, outputJSON bool) (string, error) {
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("ファイル作成エラー: %w", err)
	}
	defer file.Close()

	if !outputJSON {
		// テキスト形式で保存
		for i, review := range reviews {
			fmt.Fprintf(file, "=== レビュー %d ===\n", i+1)
			fmt.Fprintf(file, "ID: %s\n", review.RecommendationID)
			fmt.Fprintf(file, "language: %s\n", review.Language)
			fmt.Fprintf(file, "voted_up: ")
			if review.VotedUp {
				fmt.Fprintf(file, "true\n")
			} else {
				fmt.Fprintf(file, "false\n")
			}
			fmt.Fprintf(file, "votes_up: %d\n", review.VotesUp)
			fmt.Fprintf(file, "votes_funny: %d\n", review.VotesFunny)
			fmt.Fprintf(file, "weighted_score: %.2f\n", review.WeightedScore)
			fmt.Fprintf(file, "steam_purchase: %t\n", review.SteamPurchase)
			fmt.Fprintf(file, "playtime: %d分\n", review.Author.PlaytimeAtReview)
			fmt.Fprintf(file, "created_at: %s\n", time.Unix(review.TimestampCreated, 0).Format("2006-01-02 15:04:05"))
			if review.TimestampUpdated > 0 {
				fmt.Fprintf(file, "updated_at: %s\n", time.Unix(review.TimestampUpdated, 0).Format("2006-01-02 15:04:05"))
			}
			fmt.Fprintf(file, "review:\n%s\n", review.Review)
			if review.DeveloperResponse != "" {
				fmt.Fprintf(file, "developer_response:\n%s\n", review.DeveloperResponse)
				if review.TimestampDevResponse > 0 {
					fmt.Fprintf(file, "developer_response_timestamp: %s\n", time.Unix(review.TimestampDevResponse, 0).Format("2006-01-02 15:04:05"))
				}
			}
			fmt.Fprintf(file, "\n")
		}
	} else {
		// JSON形式で保存
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(reviews); err != nil {
			return "", fmt.Errorf("JSON書き込みエラー: %w", err)
		}
	}

	return filename, nil
}

// SaveReviewsByLanguage レビューを言語別に分けてファイルに保存
func SaveReviewsByLanguage(reviews []models.ReviewData, baseFilename, outputDir string, verbose bool, outputJSON bool) ([]string, error) {
	var savedFiles []string
	// 言語別にレビューを分類
	reviewsByLanguage := make(map[string][]models.ReviewData)

	for _, review := range reviews {
		lang := review.Language
		if lang == "" {
			lang = "unknown"
		}
		reviewsByLanguage[lang] = append(reviewsByLanguage[lang], review)
	}

	// ファイル拡張子を決定
	ext := ".txt"
	if outputJSON {
		ext = ".json"
	}

	// 言語別にファイル保存
	for lang, langReviews := range reviewsByLanguage {
		filename := strings.TrimSuffix(baseFilename, config.FileExtJSON) + "_" + lang + ext
		if outputDir != "" {
			filename = outputDir + "/" + filename
		}

		if savedFile, err := SaveReviewsToFile(langReviews, filename, outputJSON); err != nil {
			log.Printf("言語 %s のファイル保存エラー: %v", lang, err)
			continue
		} else {
			savedFiles = append(savedFiles, savedFile)
			if verbose {
				log.Printf("言語 %s: %d件のレビューを %s に保存", lang, len(langReviews), filename)
			}
		}
	}

	// 全体のサマリーも保存
	summaryFilename := strings.TrimSuffix(baseFilename, config.FileExtJSON) + "_all_languages" + ext
	if outputDir != "" {
		summaryFilename = outputDir + "/" + summaryFilename
	}

	if savedFile, err := SaveReviewsToFile(reviews, summaryFilename, outputJSON); err != nil {
		return nil, fmt.Errorf("サマリーファイル保存エラー: %w", err)
	} else {
		savedFiles = append(savedFiles, savedFile)
		if verbose {
			log.Printf("全言語統合ファイルを保存: %s (%d件)", summaryFilename, len(reviews))
		}
	}

	return savedFiles, nil
}
