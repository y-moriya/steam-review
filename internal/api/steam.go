package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/y-moriya/steam-review/internal/models"
	"github.com/y-moriya/steam-review/pkg/config"
)

// GetAppIDByName ゲーム名からSteam App IDを取得
func GetAppIDByName(gameName string) (string, error) {
	url := "https://api.steampowered.com/ISteamApps/GetAppList/v2/"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Steam API取得エラー: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Applist struct {
			Apps []struct {
				AppID int    `json:"appid"`
				Name  string `json:"name"`
			} `json:"apps"`
		} `json:"applist"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("JSONデコードエラー: %w", err)
	}

	for _, app := range result.Applist.Apps {
		if strings.EqualFold(app.Name, gameName) {
			return fmt.Sprintf("%d", app.AppID), nil
		}
	}
	return "", fmt.Errorf("ゲーム '%s' が見つかりません", gameName)
}

// setLanguageFilter Steam APIのリクエストパラメータに言語フィルタを設定
func setLanguageFilter(params url.Values, languages []string) {
	if len(languages) > 0 {
		// "all"が含まれている場合は全言語を取得
		for _, lang := range languages {
			if strings.ToLower(lang) == "all" {
				params.Set("language", "all")
				return
			}
		}
		// それ以外の場合は指定された言語をカンマ区切りで設定
		params.Set("language", strings.Join(languages, ","))
	} else {
		params.Set("language", "all")
	}
}

// setFilter Steam APIのリクエストパラメータにフィルターと関連パラメータを設定
func setFilter(params url.Values, filter string) {
	switch filter {
	case config.FilterRecent:
		params.Set("filter", "recent")
		params.Set("day_range", "0") // recentの場合はday_rangeは影響しない
	case config.FilterUpdated:
		params.Set("filter", "updated")
		params.Set("day_range", "0") // updatedの場合はday_rangeは影響しない
	default:
		params.Set("filter", "all")    // デフォルトは有用性による並び替え
		params.Set("day_range", "365") // allフィルターの場合、365日（最大値）を設定
	}
}

// FetchReviewsFromSteam Steam APIから直接レビューを取得
func FetchReviewsFromSteam(appID string, cursor string, numPerPage int, filter string, languages []string) (*models.SteamReviewResponse, error) {
	baseURL := "https://store.steampowered.com/appreviews/" + appID

	params := url.Values{}
	params.Set("json", "1")
	params.Set("cursor", cursor)
	params.Set("num_per_page", strconv.Itoa(numPerPage))
	params.Set("review_type", "all")
	params.Set("purchase_type", "all")

	setLanguageFilter(params, languages)
	setFilter(params, filter)

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP エラー: %d", resp.StatusCode)
	}

	var result models.SteamReviewResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON デコードエラー: %w", err)
	}

	if result.Success != 1 {
		return nil, fmt.Errorf("Steam API エラー: success = %d", result.Success)
	}

	return &result, nil
}

// FilterReviewsByLanguage 指定された言語のレビューのみをフィルタ
func FilterReviewsByLanguage(reviews []models.ReviewData, languages []string) []models.ReviewData {
	if len(languages) == 0 {
		return reviews
	}

	// "all" が指定されている場合はすべてのレビューを返す
	for _, lang := range languages {
		if strings.ToLower(lang) == "all" {
			return reviews
		}
	}

	langSet := make(map[string]bool)
	for _, lang := range languages {
		langSet[strings.ToLower(lang)] = true
	}

	var filtered []models.ReviewData
	for _, review := range reviews {
		if langSet[strings.ToLower(review.Language)] {
			filtered = append(filtered, review)
		}
	}

	return filtered
}

// FetchAllReviews 指定されたApp IDのレビューを取得
func FetchAllReviews(appID string, maxReviews int, verbose bool, languages []string, filter string) ([]models.ReviewData, error) {
	var allReviews []models.ReviewData
	cursor := "*"
	numPerPage := 100

	// 言語フィルタの準備
	langSet := make(map[string]bool)
	checkLanguage := len(languages) > 0
	if checkLanguage {
		for _, lang := range languages {
			if strings.ToLower(lang) == "all" {
				checkLanguage = false
				break
			}
			langSet[strings.ToLower(lang)] = true
		}
	}

	if verbose {
		log.Printf("App ID %s のレビュー取得を開始します", appID)
	}

	for {
		if verbose {
			log.Printf("現在のレビュー数: %d, カーソル: %s", len(allReviews), cursor)
		}

		resp, err := FetchReviewsFromSteam(appID, cursor, numPerPage, filter, languages)
		if err != nil {
			return nil, fmt.Errorf("レビュー取得エラー: %w", err)
		}

		if len(resp.Reviews) == 0 {
			if verbose {
				log.Println("これ以上レビューがありません")
			}
			break
		}

		for _, sr := range resp.Reviews {
			// 言語フィルタの適用
			if checkLanguage && !langSet[strings.ToLower(sr.Language)] {
				continue
			}

			rd := models.ConvertSteamReview(sr)
			allReviews = append(allReviews, rd)

			if maxReviews > 0 && len(allReviews) >= maxReviews {
				if verbose {
					log.Printf("最大レビュー数 %d に到達しました", maxReviews)
				}
				return allReviews[:maxReviews], nil
			}
		}

		if resp.Cursor == cursor || resp.Cursor == "" {
			if verbose {
				log.Println("カーソルが変更されませんでした。終了します")
			}
			break
		}

		cursor = resp.Cursor

		// レート制限対策
		time.Sleep(1 * time.Second)
	}

	if verbose {
		log.Printf("合計 %d 件のレビューを取得しました", len(allReviews))
	}
	return allReviews, nil
}

// GetReviewsByGameName ゲーム名からレビューを取得
func GetReviewsByGameName(gameName string, maxReviews int, verbose bool, languages []string, filter string) ([]models.ReviewData, string, error) {
	appID, err := GetAppIDByName(gameName)
	if err != nil {
		return nil, "", fmt.Errorf("App ID取得エラー: %w", err)
	}
	if verbose {
		log.Printf("ゲーム '%s' (App ID: %s) のレビューを取得します", gameName, appID)
	}
	reviews, err := FetchAllReviews(appID, maxReviews, verbose, languages, filter)
	return reviews, appID, err
}
