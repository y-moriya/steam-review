package main

import (
	"encoding/json"
	"flag"
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

const (
	// バージョン情報
	Version = "v0.3.1"     // プログラムのバージョン
	AppName = "Steam Reviews CLI Tool" // プログラム名

	// レビューのフィルター
	FilterAll     = "all"     // 有用性による並び替え
	FilterRecent  = "recent"  // 作成日時による並び替え
	FilterUpdated = "updated" // 最終更新日時による並び替え

	// ファイル形式
	FileExtJSON = ".json" // JSON形式のファイル拡張子
	FileExtTXT  = ".txt"  // テキスト形式のファイル拡張子
)

	// コマンドライン引数の設定
type Config struct {
	AppID        string
	GameName     string
	MaxReviews   int
	Languages    []string
	OutputDir    string
	Verbose      bool
	SplitByLang  bool
	OutputJSON   bool
	Filter      string   // レビューのフィルター
}// GetAppIDByName ゲーム名からSteam App IDを取得
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
	case FilterRecent:
		params.Set("filter", "recent")
		params.Set("day_range", "0") // recentの場合はday_rangeは影響しない
	case FilterUpdated:
		params.Set("filter", "updated")
		params.Set("day_range", "0") // updatedの場合はday_rangeは影響しない
	default:
		params.Set("filter", "all") // デフォルトは有用性による並び替え
		params.Set("day_range", "365") // allフィルターの場合、365日（最大値）を設定
	}
}

// FetchReviewsFromSteam Steam APIから直接レビューを取得
func FetchReviewsFromSteam(appID string, cursor string, numPerPage int, filter string, languages []string) (*SteamReviewResponse, error) {
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
		WeightedScore:        float64(sr.WeightedVoteScore),
		CommentCount:         sr.CommentCount,
		SteamPurchase:        sr.SteamPurchase,
		ReceivedForFree:      sr.ReceivedForFree,
		WrittenDuringEA:      sr.WrittenDuringEA,
		DeveloperResponse:    sr.DeveloperResponse,
		TimestampDevResponse: sr.TimestampDevResp,
	}
}

	// FilterReviewsByLanguage 指定された言語のレビューのみをフィルタ
func FilterReviewsByLanguage(reviews []ReviewData, languages []string) []ReviewData {
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
	
	var filtered []ReviewData
	for _, review := range reviews {
		if langSet[strings.ToLower(review.Language)] {
			filtered = append(filtered, review)
		}
	}
	
	return filtered
}

// FetchAllReviews 指定されたApp IDのレビューを取得
func FetchAllReviews(appID string, maxReviews int, verbose bool, languages []string, filter string) ([]ReviewData, error) {
	var allReviews []ReviewData
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

			rd := ConvertSteamReview(sr)
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

// SaveReviewsToFile レビューをファイルに保存
func SaveReviewsToFile(reviews []ReviewData, filename string, outputJSON bool) (string, error) {
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
			fmt.Fprintf(file, "言語: %s\n", review.Language)
			fmt.Fprintf(file, "評価: ")
			if review.VotedUp {
				fmt.Fprintf(file, "肯定的\n")
			} else {
				fmt.Fprintf(file, "否定的\n")
			}
			fmt.Fprintf(file, "投票数: %d\n", review.VotesUp)
			fmt.Fprintf(file, "面白い投票: %d\n", review.VotesFunny)
			fmt.Fprintf(file, "重み付けスコア: %.2f\n", review.WeightedScore)
			fmt.Fprintf(file, "Steam購入: %t\n", review.SteamPurchase)
			fmt.Fprintf(file, "プレイ時間: %d分\n", review.Author.PlaytimeAtReview)
			fmt.Fprintf(file, "作成日時: %s\n", time.Unix(review.TimestampCreated, 0).Format("2006-01-02 15:04:05"))
			if review.TimestampUpdated > 0 {
				fmt.Fprintf(file, "更新日時: %s\n", time.Unix(review.TimestampUpdated, 0).Format("2006-01-02 15:04:05"))
			}
			fmt.Fprintf(file, "レビュー内容:\n%s\n", review.Review)
			if review.DeveloperResponse != "" {
				fmt.Fprintf(file, "\n開発者応答:\n%s\n", review.DeveloperResponse)
				if review.TimestampDevResponse > 0 {
					fmt.Fprintf(file, "開発者応答日時: %s\n", time.Unix(review.TimestampDevResponse, 0).Format("2006-01-02 15:04:05"))
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
func SaveReviewsByLanguage(reviews []ReviewData, baseFilename, outputDir string, verbose bool, outputJSON bool) ([]string, error) {
	var savedFiles []string
	// 言語別にレビューを分類
	reviewsByLanguage := make(map[string][]ReviewData)
	
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
		filename := strings.TrimSuffix(baseFilename, FileExtJSON) + "_" + lang + ext
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
	summaryFilename := strings.TrimSuffix(baseFilename, FileExtJSON) + "_all_languages" + ext
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

// GetReviewsByGameName ゲーム名からレビューを取得
func GetReviewsByGameName(gameName string, maxReviews int, verbose bool, languages []string, filter string) ([]ReviewData, string, error) {
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

// parseLanguages カンマ区切りの言語文字列を配列に変換
func parseLanguages(langStr string) []string {
	var languages []string
	for _, lang := range strings.Split(langStr, ",") {
		lang = strings.TrimSpace(lang)
		if lang != "" {
			languages = append(languages, lang)
		}
	}
	return languages
}

// printUsage 使用方法を表示
func printUsage() {
	fmt.Printf(`%s version %s

使用方法:
  %s [オプション]

オプション:
  -appid string         Steam App ID (例: 440)
  -game string          ゲーム名 (例: "Team Fortress 2")
  -max int             最大取得レビュー数 (デフォルト: 100, 0で無制限)
  -lang string         取得する言語 (カンマ区切り, デフォルト: japanese, 例: "japanese,english")
  -output string       出力ディレクトリ (デフォルト: output)
  -split              言語別にファイルを分けて保存
  -json               出力ファイルをJSON形式(.json)にする (デフォルト: テキスト形式)
  -verbose            詳細なログを表示
  -filter string      レビューのフィルター (recent: 作成日時順, updated: 更新日時順, all: 有用性順(デフォルト))
  -help               このヘルプを表示
  -version            バージョン情報を表示

使用例:
  # App IDを指定して日本語レビューを取得（デフォルト: 有用性順）
  %s -appid 440 -max 500 -verbose

  # 作成日時順でレビューを取得
  %s -appid 440 -max 500 -filter recent -verbose

  # ゲーム名で英語レビューを取得
  %s -game "Cyberpunk 2077" -lang "english" -max 1000 -output ./reviews

  # 複数言語のレビューを取得
  %s -game "Elden Ring" -lang "japanese,english" -max 300 -split

  # 日本語レビューをJSON形式で保存
  %s -appid 570 -max 2000 -output ./dota2_reviews -json -verbose

  # すべての言語のレビューを取得
  %s -appid 730 -lang "all" -max 1000 -split

注意:
  - App IDとゲーム名のどちらか一方を指定してください
  - -lang を指定しない場合、デフォルトで日本語レビューのみを取得します
  - "all" を指定するとすべての言語のレビューを取得します
  - 大量のレビューを取得する場合は時間がかかります
  - Steam APIのレート制限により、リクエスト間に1秒の待機時間があります
`, AppName, Version, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	var config Config
	var languageStr string
	var help bool
	var showVersion bool

	// コマンドライン引数の定義
	flag.BoolVar(&showVersion, "version", false, "バージョン情報を表示")
	flag.StringVar(&config.AppID, "appid", "", "Steam App ID")
	flag.StringVar(&config.GameName, "game", "", "ゲーム名")
	flag.IntVar(&config.MaxReviews, "max", 100, "最大取得レビュー数 (0で無制限)")
	flag.StringVar(&languageStr, "lang", "japanese", "取得する言語 (カンマ区切り, デフォルト: japanese)")
	flag.StringVar(&config.OutputDir, "output", "output", "出力ディレクトリ")
	flag.StringVar(&config.Filter, "filter", FilterAll, 
		"レビューのフィルター (recent: 作成日時順, updated: 更新日時順, all: 有用性順(デフォルト))")
	flag.BoolVar(&config.SplitByLang, "split", false, "言語別にファイルを分けて保存")
	flag.BoolVar(&config.OutputJSON, "json", false, "出力ファイルをJSON形式(.json)にする (デフォルト: テキスト形式)")
	flag.BoolVar(&config.Verbose, "verbose", false, "詳細なログを表示")
	flag.BoolVar(&help, "help", false, "ヘルプを表示")

	flag.Parse()

	// バージョン情報の表示
	if showVersion {
		fmt.Printf("%s version %s\n", AppName, Version)
		return
	}

	// ヘルプ表示
	if help {
		printUsage()
		return
	}

	// 言語設定をパース
	config.Languages = parseLanguages(languageStr)

	// バリデーション
	if config.AppID == "" && config.GameName == "" {
		fmt.Printf("エラー: App ID またはゲーム名を指定してください\n\n")
		printUsage()
		os.Exit(1)
	}

	if config.AppID != "" && config.GameName != "" {
		fmt.Printf("エラー: App ID とゲーム名の両方を指定することはできません\n\n")
		printUsage()
		os.Exit(1)
	}

	// 出力ディレクトリの作成
	if config.OutputDir != "" {
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			log.Fatalf("出力ディレクトリの作成に失敗しました: %v", err)
		}
	}

	var reviews []ReviewData
	var appID string
	var gameName string
	var err error

	// レビュー取得
	if config.AppID != "" {
		appID = config.AppID
		gameName = fmt.Sprintf("App ID %s", appID)
		if config.Verbose {
			log.Printf("App ID %s のレビューを取得中...", appID)
		}
		reviews, err = FetchAllReviews(appID, config.MaxReviews, config.Verbose, config.Languages, config.Filter)
	} else {
		gameName = config.GameName
		if config.Verbose {
			log.Printf("ゲーム '%s' のレビューを取得中...", gameName)
		}
		reviews, appID, err = GetReviewsByGameName(gameName, config.MaxReviews, config.Verbose, config.Languages, config.Filter)
	}

	if err != nil {
		log.Fatalf("レビュー取得エラー: %v", err)
	}

	if len(reviews) == 0 {
		log.Printf("レビューが見つかりませんでした")
		return
	}

	// ファイル保存
	ext := ".txt"
	if config.OutputJSON {
		ext = ".json"
	}
	baseFilename := fmt.Sprintf("steam_reviews_%s%s", appID, ext)
	
	var savedFiles []string
	if config.SplitByLang {
		files, err := SaveReviewsByLanguage(reviews, baseFilename, config.OutputDir, config.Verbose, config.OutputJSON)
		if err != nil {
			log.Printf("ファイル保存エラー: %v", err)
		}
		savedFiles = files
	} else {
		filename := baseFilename
		if config.OutputDir != "" {
			filename = config.OutputDir + "/" + filename
		}
		
		if savedFile, err := SaveReviewsToFile(reviews, filename, config.OutputJSON); err != nil {
			log.Printf("ファイル保存エラー: %v", err)
		} else {
			savedFiles = append(savedFiles, savedFile)
			if config.Verbose {
				log.Printf("レビューを %s に保存しました", filename)
			}
		}
	}

	// 保存したファイル一覧を表示
	fmt.Printf("\n=== 保存したファイル一覧 ===\n")
	for _, file := range savedFiles {
		fmt.Printf("- %s\n", file)
	}
	fmt.Println()

	// 統計情報を表示
	PrintReviewStats(reviews, gameName)
}