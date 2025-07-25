package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// GameDetails ゲーム詳細情報構造体
type GameDetails struct {
	AppID       string    `json:"app_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Publisher   []string  `json:"publisher"`
	Developer   []string  `json:"developer"`
	ReleaseDate string    `json:"release_date"`
	Price       string    `json:"price"`
	Currency    string    `json:"currency"`
	Tags        []string  `json:"tags"`
	Categories  []string  `json:"categories"`
	Genres      []string  `json:"genres"`
	HeaderImage string    `json:"header_image"`
	Website     string    `json:"website"`
	RequiredAge int       `json:"required_age"`
	IsFree      bool      `json:"is_free"`
	RetrievedAt time.Time `json:"retrieved_at"`
}

// SteamAppDetailsResponse Steam Store APIからのレスポンス構造体
type SteamAppDetailsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Name                string   `json:"name"`
		SteamAppID          int      `json:"steam_appid"`
		RequiredAge         int      `json:"required_age"`
		IsFree              bool     `json:"is_free"`
		DetailedDescription string   `json:"detailed_description"`
		AboutTheGame        string   `json:"about_the_game"`
		ShortDescription    string   `json:"short_description"`
		HeaderImage         string   `json:"header_image"`
		Website             string   `json:"website"`
		Developers          []string `json:"developers"`
		Publishers          []string `json:"publishers"`
		PriceOverview       struct {
			Currency        string `json:"currency"`
			Initial         int    `json:"initial"`
			Final           int    `json:"final"`
			DiscountPercent int    `json:"discount_percent"`
			FinalFormatted  string `json:"final_formatted"`
		} `json:"price_overview"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
	} `json:"data"`
}

// ConvertToGameDetails Steam Store APIのレスポンスをGameDetailsに変換
func ConvertToGameDetails(appID string, response SteamAppDetailsResponse) GameDetails {
	var categories []string
	for _, cat := range response.Data.Categories {
		categories = append(categories, cat.Description)
	}

	var genres []string
	for _, genre := range response.Data.Genres {
		genres = append(genres, genre.Description)
	}

	price := "Free"
	currency := ""
	if !response.Data.IsFree && response.Data.PriceOverview.FinalFormatted != "" {
		price = response.Data.PriceOverview.FinalFormatted
		currency = response.Data.PriceOverview.Currency
	}

	// 説明文は短い説明を優先し、なければ詳細説明を使用
	description := response.Data.ShortDescription
	if description == "" {
		description = response.Data.AboutTheGame
		if description == "" {
			description = response.Data.DetailedDescription
		}
	}

	// HTMLタグを簡易的に除去
	description = strings.ReplaceAll(description, "<br>", "\n")
	description = strings.ReplaceAll(description, "<p>", "")
	description = strings.ReplaceAll(description, "</p>", "\n")

	return GameDetails{
		AppID:       appID,
		Name:        response.Data.Name,
		Description: description,
		Publisher:   response.Data.Publishers,
		Developer:   response.Data.Developers,
		ReleaseDate: response.Data.ReleaseDate.Date,
		Price:       price,
		Currency:    currency,
		Categories:  categories,
		Genres:      genres,
		HeaderImage: response.Data.HeaderImage,
		Website:     response.Data.Website,
		RequiredAge: response.Data.RequiredAge,
		IsFree:      response.Data.IsFree,
		RetrievedAt: time.Now(),
	}
}

// SteamReviewResponse Steam APIからのレスポンス構造体
type SteamReviewResponse struct {
	Success      int `json:"success"`
	QuerySummary struct {
		NumReviews      int    `json:"num_reviews"`
		ReviewScore     int    `json:"review_score"`
		ReviewScoreDesc string `json:"review_score_desc"`
		TotalPositive   int    `json:"total_positive"`
		TotalNegative   int    `json:"total_negative"`
		TotalReviews    int    `json:"total_reviews"`
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
	RecommendationID  string          `json:"recommendationid"`
	Author            SteamAuthor     `json:"author"`
	Language          string          `json:"language"`
	Review            string          `json:"review"`
	TimestampCreated  int64           `json:"timestamp_created"`
	TimestampUpdated  int64           `json:"timestamp_updated"`
	VotedUp           bool            `json:"voted_up"`
	VotesUp           int             `json:"votes_up"`
	VotesFunny        int             `json:"votes_funny"`
	WeightedVoteScore FlexibleFloat64 `json:"weighted_vote_score"`
	CommentCount      int             `json:"comment_count"`
	SteamPurchase     bool            `json:"steam_purchase"`
	ReceivedForFree   bool            `json:"received_for_free"`
	WrittenDuringEA   bool            `json:"written_during_early_access"`
	DeveloperResponse string          `json:"developer_response"`
	TimestampDevResp  int64           `json:"timestamp_dev_responded"`
}

// SteamAuthor レビュー作者の構造体
type SteamAuthor struct {
	SteamID              string `json:"steamid"`
	NumGamesOwned        int    `json:"num_games_owned"`
	NumReviews           int    `json:"num_reviews"`
	PlayTimeForever      int    `json:"playtime_forever"`
	PlayTimeLastTwoWeeks int    `json:"playtime_last_two_weeks"`
	PlayTimeAtReview     int    `json:"playtime_at_review"`
	LastPlayed           int64  `json:"last_played"`
}

// ReviewData 最終的なレビューデータ構造体
type ReviewData struct {
	RecommendationID     string     `json:"recommendation_id"`
	Author               AuthorData `json:"author"`
	Language             string     `json:"language"`
	Review               string     `json:"review"`
	TimestampCreated     int64      `json:"timestamp_created"`
	TimestampUpdated     int64      `json:"timestamp_updated"`
	VotedUp              bool       `json:"voted_up"`
	VotesUp              int        `json:"votes_up"`
	VotesFunny           int        `json:"votes_funny"`
	WeightedScore        float64    `json:"weighted_vote_score"`
	CommentCount         int        `json:"comment_count"`
	SteamPurchase        bool       `json:"steam_purchase"`
	ReceivedForFree      bool       `json:"received_for_free"`
	WrittenDuringEA      bool       `json:"written_during_early_access"`
	DeveloperResponse    string     `json:"developer_response,omitempty"`
	TimestampDevResponse int64      `json:"timestamp_dev_responded,omitempty"`
}

// AuthorData 作成者データ構造体
type AuthorData struct {
	SteamID              string `json:"steam_id"`
	NumGamesOwned        int    `json:"num_games_owned"`
	NumReviews           int    `json:"num_reviews"`
	PlaytimeForever      int    `json:"playtime_forever"`
	PlaytimeLastTwoWeeks int    `json:"playtime_last_two_weeks"`
	PlaytimeAtReview     int    `json:"playtime_at_review"`
	LastPlayed           int64  `json:"last_played"`
}

// ConvertSteamReview Steam APIのレビューを内部形式に変換
func ConvertSteamReview(sr SteamReview) ReviewData {
	return ReviewData{
		RecommendationID: sr.RecommendationID,
		Author: AuthorData{
			SteamID:              sr.Author.SteamID,
			NumGamesOwned:        sr.Author.NumGamesOwned,
			NumReviews:           sr.Author.NumReviews,
			PlaytimeForever:      sr.Author.PlayTimeForever,
			PlaytimeLastTwoWeeks: sr.Author.PlayTimeLastTwoWeeks,
			PlaytimeAtReview:     sr.Author.PlayTimeAtReview,
			LastPlayed:           sr.Author.LastPlayed,
		},
		Language:             sr.Language,
		Review:               sr.Review,
		TimestampCreated:     sr.TimestampCreated,
		TimestampUpdated:     sr.TimestampUpdated,
		VotedUp:              sr.VotedUp,
		VotesUp:              sr.VotesUp,
		VotesFunny:           sr.VotesFunny,
		WeightedScore:        float64(sr.WeightedVoteScore),
		CommentCount:         sr.CommentCount,
		SteamPurchase:        sr.SteamPurchase,
		ReceivedForFree:      sr.ReceivedForFree,
		WrittenDuringEA:      sr.WrittenDuringEA,
		DeveloperResponse:    sr.DeveloperResponse,
		TimestampDevResponse: sr.TimestampDevResp,
	}
}
