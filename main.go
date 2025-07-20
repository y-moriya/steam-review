package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// SteamReviewResponse Steam APIからのレスポンス構造体
type SteamReviewResponse struct {
	Success      int    `json:"success"`
	QuerySummary struct {
		NumReviews      int     `json:"num_reviews"`
		ReviewScore     int     `json:"review_score"`
		ReviewScoreDesc string  `json:"review_score_desc"`
		TotalPositive   int     `json:"total_positive"`
		TotalNegative   int     `json:"total_negative"`
		TotalReviews    int     `json:"total_reviews"`
	} `json:"query_summary"`
	Reviews []SteamReview `json:"reviews"`
	Cursor  string        `json:"cursor"`
}

// FlexibleFloat64 文字列または数値を float64 として受け取るカスタム型
type FlexibleFloat64 float64

func (f *FlexibleFloat64) UnmarshalJSON(data []byte) error {
	// 数値として試す
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexibleFloat64(num)
		return nil
	}
	
	// 文字列として試す
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" {
			*f = FlexibleFloat64(0)
			return nil
		}
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			*f = FlexibleFloat64(num)
			return nil
		}
	}
	
	// デフォルト値を設定
	*f = FlexibleFloat64(0)
	return nil
}

// SteamReview 個別のレビュー構造体
type SteamReview struct {
	RecommendationID     string          `json:"recommendationid"`
	Author              SteamAuthor     `json:"author"`
	Language            string          `json:"language"`
	Review              string          `json:"review"`
	TimestampCreated    int64           `json:"timestamp_created"`
	TimestampUpdated    int64           `json:"timestamp_updated"`
	VotedUp             bool            `json:"voted_up"`
	VotesUp             int             `json:"votes_up"`
	VotesFunny          int             `json:"votes_funny"`
	WeightedVoteScore   FlexibleFloat64 `json:"weighted_vote_score"`
	CommentCount        int             `json:"comment_count"`
	SteamPurchase       bool            `json:"steam_purchase"`
	ReceivedForFree     bool            `json:"received_for_free"`
	WrittenDuringEA     bool            `json:"written_during_early_access"`
	DeveloperResponse   string          `json:"developer_response"`
	TimestampDevResp    int64           `json:"timestamp_dev_responded"`
}

// SteamAuthor レビュー作者の構造体
type SteamAuthor struct {
	SteamID                string `json:"steamid"`
	NumGamesOwned         int    `json:"num_games_owned"`
	NumReviews            int    `json:"num_reviews"`
	PlayTimeForever       int    `json:"playtime_forever"`
	PlayTimeLastTwoWeeks  int    `json:"playtime_last_two_weeks"`
	PlayTimeAtReview      int    `json:"playtime_at_review"`
	LastPlayed            int64  `json:"last_played"`
}

// ReviewData 最終的なレビューデータ構造体
type ReviewData struct {
	RecommendationID      string      `json:"recommendation_id"`
	Author               AuthorData   `json:"author"`
	Language             string       `json:"language"`
	Review               string       `json:"review"`
	TimestampCreated     int64        `json:"timestamp_created"`
	TimestampUpdated     int64        `json:"timestamp_updated"`
	VotedUp              bool         `json:"voted_up"`
	VotesUp              int          `json:"votes_up"`
	VotesFunny           int          `json:"votes_funny"`
	WeightedScore        float64      `json:"weighted_vote_score"`
	CommentCount         int          `json:"comment_count"`
	SteamPurchase        bool         `json:"steam_purchase"`
	ReceivedForFree      bool         `json:"received_for_free"`
	WrittenDuringEA      bool         `json:"written_during_early_access"`
	DeveloperResponse    string       `json:"developer_response,omitempty"`
	TimestampDevResponse int64        `json:"timestamp_dev_responded,omitempty"`
}

// AuthorData 作成者データ構造体
type AuthorData struct {
	SteamID                 string `json:"steam_id"`
	NumGamesOwned          int    `json:"num_games_owned"`
	NumReviews             int    `json:"num_reviews"`
	PlaytimeForever        int    `json:"playtime_forever"`
	PlaytimeLastTwoWeeks   int    `json:"playtime_last_two_weeks"`
	PlaytimeAtReview       int    `json:"playtime_at_review"`
	LastPlayed             int64  `json:"last_played"`
}

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

// FetchReviewsFromSteam Steam APIから直接レビューを取得
func FetchReviewsFromSteam(appID string, cursor string, numPerPage int) (*SteamReviewResponse, error) {
	baseURL := "https://store.steampowered.com/appreviews/" + appID
	
	params := url.Values{}
	params.Set("json", "1")
	params.Set("cursor", cursor)
	params.Set("num_per_page", strconv.Itoa(numPerPage))
	params.Set("filter", "all")
	params.Set("language", "all")
	params.Set("day_range", "9223372036854775807")
	params.Set("review_type", "all")
	params.Set("purchase_type", "all")
	
	fullURL := baseURL + "?" + params.Encode()
	
	log.Printf("リクエストURL: %s", fullURL)
	
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストエラー: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP エラー: %d", resp.StatusCode)
	}
	
	var result SteamReviewResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON デコードエラー: %w", err)
	}
	
	if result.Success != 1 {
		return nil, fmt.Errorf("Steam API エラー: success = %d", result.Success)
	}
	
	return &result, nil
}

// ConvertSteamReview Steam APIのレビューを内部形式に変換
func ConvertSteamReview(sr SteamReview) ReviewData {
	return ReviewData{
		RecommendationID: sr.RecommendationID,
		Author: AuthorData{
			SteamID:                sr.Author.SteamID,
			NumGamesOwned:         sr.Author.NumGamesOwned,
			NumReviews:            sr.Author.NumReviews,
			PlaytimeForever:       sr.Author.PlayTimeForever,
			PlaytimeLastTwoWeeks:  sr.Author.PlayTimeLastTwoWeeks,
			PlaytimeAtReview:      sr.Author.PlayTimeAtReview,
			LastPlayed:            sr.Author.LastPlayed,
		},
		Language:             sr.Language,
		Review:               sr.Review,
		TimestampCreated:     sr.TimestampCreated,
		TimestampUpdated:     sr.TimestampUpdated,
		VotedUp:              sr.VotedUp,
		VotesUp:              sr.VotesUp,
		VotesFunny:           sr.VotesFunny,
		WeightedScore:        float64(sr.WeightedVoteScore), // FlexibleFloat64 から float64 に変換
		CommentCount:         sr.CommentCount,
		SteamPurchase:        sr.SteamPurchase,
		ReceivedForFree:      sr.ReceivedForFree,
		WrittenDuringEA:      sr.WrittenDuringEA,
		DeveloperResponse:    sr.DeveloperResponse,
		TimestampDevResponse: sr.TimestampDevResp,
	}
}

// FetchAllReviews 指定されたApp IDのレビューを取得
func FetchAllReviews(appID string, maxReviews int) ([]ReviewData, error) {
	var allReviews []ReviewData
	cursor := "*"
	numPerPage := 100
	
	log.Printf("App ID %s のレビュー取得を開始します", appID)
	
	for {
		log.Printf("現在のレビュー数: %d, カーソル: %s", len(allReviews), cursor)
		
		resp, err := FetchReviewsFromSteam(appID, cursor, numPerPage)
		if err != nil {
			return nil, fmt.Errorf("レビュー取得エラー: %w", err)
		}
		
		if len(resp.Reviews) == 0 {
			log.Println("これ以上レビューがありません")
			break
		}
		
		for _, sr := range resp.Reviews {
			rd := ConvertSteamReview(sr)
			allReviews = append(allReviews, rd)
			
			if maxReviews > 0 && len(allReviews) >= maxReviews {
				log.Printf("最大レビュー数 %d に到達しました", maxReviews)
				return allReviews, nil
			}
		}
		
		if resp.Cursor == cursor || resp.Cursor == "" {
			log.Println("カーソルが変更されませんでした。終了します")
			break
		}
		
		cursor = resp.Cursor
		
		// レート制限対策
		time.Sleep(1 * time.Second)
	}
	
	log.Printf("合計 %d 件のレビューを取得しました", len(allReviews))
	return allReviews, nil
}

// SaveReviewsToFile レビューをJSONファイルに保存
func SaveReviewsToFile(reviews []ReviewData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ファイル作成エラー: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(reviews); err != nil {
		return fmt.Errorf("JSON書き込みエラー: %w", err)
	}

	log.Printf("レビューを %s に保存しました", filename)
	return nil
}

// SaveReviewsByLanguage レビューを言語別に分けてJSONファイルに保存
func SaveReviewsByLanguage(reviews []ReviewData, baseFilename string) error {
	// 言語別にレビューを分類
	reviewsByLanguage := make(map[string][]ReviewData)
	
	for _, review := range reviews {
		lang := review.Language
		if lang == "" {
			lang = "unknown"
		}
		reviewsByLanguage[lang] = append(reviewsByLanguage[lang], review)
	}
	
	// 言語別にファイル保存
	for lang, langReviews := range reviewsByLanguage {
		// ファイル名を生成（拡張子を除去してから言語コードを追加）
		filename := strings.TrimSuffix(baseFilename, ".json") + "_" + lang + ".json"
		
		if err := SaveReviewsToFile(langReviews, filename); err != nil {
			log.Printf("言語 %s のファイル保存エラー: %v", lang, err)
			continue
		}
		
		log.Printf("言語 %s: %d件のレビューを保存", lang, len(langReviews))
	}
	
	// 全体のサマリーも保存
	summaryFilename := strings.TrimSuffix(baseFilename, ".json") + "_all_languages.json"
	if err := SaveReviewsToFile(reviews, summaryFilename); err != nil {
		return fmt.Errorf("サマリーファイル保存エラー: %w", err)
	}
	
	log.Printf("全言語統合ファイルを保存: %s (%d件)", summaryFilename, len(reviews))
	
	return nil
}

// GetReviewsByGameName ゲーム名からレビューを取得
func GetReviewsByGameName(gameName string, maxReviews int) ([]ReviewData, error) {
	appID, err := GetAppIDByName(gameName)
	if err != nil {
		return nil, fmt.Errorf("App ID取得エラー: %w", err)
	}
	log.Printf("ゲーム '%s' (App ID: %s) のレビューを取得します", gameName, appID)
	return FetchAllReviews(appID, maxReviews)
}

// PrintReviewStats レビュー統計を表示
func PrintReviewStats(reviews []ReviewData, gameName string) {
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

func main() {
	// 使用例1: App IDを直接指定
	appID := "2089600"
	maxReviews := 5_000_000
	
	log.Printf("App ID %s のレビューを取得中...", appID)
	reviews, err := FetchAllReviews(appID, maxReviews)
	if err != nil {
		log.Printf("レビュー取得エラー: %v", err)
	} else {
		// 言語別にファイル保存
		baseFilename := fmt.Sprintf("steam_reviews_%s.json", appID)
		if err := SaveReviewsByLanguage(reviews, baseFilename); err != nil {
			log.Printf("ファイル保存エラー: %v", err)
		}
		PrintReviewStats(reviews, "都市伝説解体センター")
	}

}