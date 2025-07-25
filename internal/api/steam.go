package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/y-moriya/steam-review/internal/logger"
	"github.com/y-moriya/steam-review/internal/models"
	"github.com/y-moriya/steam-review/pkg/config"
	"github.com/y-moriya/steam-review/pkg/i18n"
)

// GetAppIDByName ゲーム名からSteam App IDを取得
func GetAppIDByName(gameName string) (string, error) {
	url := "https://api.steampowered.com/ISteamApps/GetAppList/v2/"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New(i18n.Tf(i18n.MsgErrorSteamAPIFetch, err))
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
		return "", errors.New(i18n.Tf(i18n.MsgErrorJSONDecode, err))
	}

	for _, app := range result.Applist.Apps {
		if strings.EqualFold(app.Name, gameName) {
			return fmt.Sprintf("%d", app.AppID), nil
		}
	}
	return "", errors.New(i18n.Tf(i18n.MsgErrorGameNotFound, gameName))
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
		return nil, errors.New(i18n.Tf(i18n.MsgErrorHTTPRequest, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(i18n.Tf(i18n.MsgErrorHTTPStatus, resp.StatusCode))
	}

	var result models.SteamReviewResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.New(i18n.Tf(i18n.MsgErrorJSONDecode, err))
	}

	if result.Success != 1 {
		return nil, errors.New(i18n.Tf(i18n.MsgErrorSteamAPIResponse, result.Success))
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
func FetchAllReviews(appID string, maxReviews int, verbose bool, languages []string, filter string, logger *logger.Logger) ([]models.ReviewData, error) {
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

	if verbose && logger != nil {
		logger.Verbose(i18n.Tf(i18n.MsgVerboseReviewFetchStart, appID))
	}

	for {
		if verbose && logger != nil {
			logger.Verbose(i18n.Tf(i18n.MsgVerboseReviewProgress, len(allReviews), cursor))
		}

		resp, err := FetchReviewsFromSteam(appID, cursor, numPerPage, filter, languages)
		if err != nil {
			return nil, errors.New(i18n.Tf(i18n.MsgErrorReviewFetch, err))
		}

		if len(resp.Reviews) == 0 {
			if verbose && logger != nil {
				logger.Verbose(i18n.T(i18n.MsgVerboseNoMoreReviews))
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
				if verbose && logger != nil {
					logger.Verbose(i18n.Tf(i18n.MsgVerboseMaxReviewsReached, maxReviews))
				}
				return allReviews[:maxReviews], nil
			}
		}

		if resp.Cursor == cursor || resp.Cursor == "" {
			if verbose && logger != nil {
				logger.Verbose(i18n.T(i18n.MsgVerboseCursorNotChanged))
			}
			break
		}

		cursor = resp.Cursor

		// レート制限対策
		time.Sleep(1 * time.Second)
	}

	if verbose && logger != nil {
		logger.Verbose(i18n.Tf(i18n.MsgVerboseTotalReviewsFetched, len(allReviews)))
	}
	return allReviews, nil
}

// GetReviewsByGameName ゲーム名からレビューを取得
func GetReviewsByGameName(gameName string, maxReviews int, verbose bool, languages []string, filter string, logger *logger.Logger) ([]models.ReviewData, string, error) {
	appID, err := GetAppIDByName(gameName)
	if err != nil {
		return nil, "", errors.New(i18n.Tf(i18n.MsgErrorAppIDFetch, err))
	}
	if verbose && logger != nil {
		logger.Verbose(i18n.Tf(i18n.MsgVerboseGameReviewFetch, gameName, appID))
	}
	reviews, err := FetchAllReviews(appID, maxReviews, verbose, languages, filter, logger)
	return reviews, appID, err
}

// GetGameDetails Steam Store APIからゲーム詳細情報を取得
func GetGameDetails(appID string, verbose bool, logger *logger.Logger) (*models.GameDetails, error) {
	if verbose && logger != nil {
		logger.Verbose(i18n.Tf(i18n.MsgVerboseGameDetailsFetch, appID))
	}

	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&l=japanese", appID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New(i18n.Tf(i18n.MsgErrorSteamStoreFetch, err))
	}
	defer resp.Body.Close()

	// レスポンスを一度マップとして読み込む
	var responseMap map[string]models.SteamAppDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, errors.New(i18n.Tf(i18n.MsgErrorJSONDecode, err))
	}

	// App IDに対応するデータを取得
	appResponse, exists := responseMap[appID]
	if !exists {
		return nil, errors.New(i18n.Tf(i18n.MsgErrorAppDataNotFound, appID))
	}

	if !appResponse.Success {
		return nil, errors.New(i18n.Tf(i18n.MsgErrorGameDetailsFail, appID))
	}

	// GameDetailsに変換
	gameDetails := models.ConvertToGameDetails(appID, appResponse)

	if verbose && logger != nil {
		logger.Verbose(i18n.Tf(i18n.MsgVerboseGameDetailsObtained, gameDetails.Name))
	}

	return &gameDetails, nil
}
