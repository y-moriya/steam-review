package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/y-moriya/steam-review/internal/api"
	"github.com/y-moriya/steam-review/internal/logger"
	"github.com/y-moriya/steam-review/internal/models"
	"github.com/y-moriya/steam-review/internal/stats"
	"github.com/y-moriya/steam-review/internal/storage"
	"github.com/y-moriya/steam-review/pkg/config"
	"github.com/y-moriya/steam-review/pkg/i18n"
)

// ParseLanguages カンマ区切りの言語文字列を配列に変換
func ParseLanguages(langStr string) []string {
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
	fmt.Printf(i18n.T(i18n.MsgUsageFull), config.AppName, config.Version)
}

func main() {
	// i18n システムを初期化
	i18n.Init()

	var cfg config.Config
	var languageStr string
	var help bool
	var showVersion bool

	// コマンドライン引数の定義
	flag.BoolVar(&showVersion, "version", false, "バージョン情報を表示")
	flag.StringVar(&cfg.AppID, "appid", "", "Steam App ID")
	flag.StringVar(&cfg.GameName, "game", "", "ゲーム名")
	flag.IntVar(&cfg.MaxReviews, "max", 100, "最大取得レビュー数 (0で無制限)")
	flag.StringVar(&languageStr, "lang", "japanese", "取得する言語 (カンマ区切り, デフォルト: japanese)")
	flag.StringVar(&cfg.OutputDir, "output", "output", "出力ディレクトリ")
	flag.StringVar(&cfg.Filter, "filter", config.FilterAll,
		"レビューのフィルター (recent: 作成日時順, updated: 更新日時順, all: 有用性順(デフォルト))")
	flag.BoolVar(&cfg.SplitByLang, "split", false, "言語別にファイルを分けて保存")
	flag.BoolVar(&cfg.OutputJSON, "json", false, "出力ファイルをJSON形式(.json)にする (デフォルト: テキスト形式)")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "詳細なログを表示")
	flag.BoolVar(&help, "help", false, "ヘルプを表示")

	flag.Parse()

	// バージョン情報の表示
	if showVersion {
		fmt.Printf("%s\n", i18n.Tf(i18n.MsgAppVersion, config.Version))
		return
	}

	// ヘルプ表示
	if help {
		printUsage()
		return
	}

	// ロガーを初期化
	log, err := logger.New("logs", cfg.Verbose)
	if err != nil {
		fmt.Printf("%s\n", i18n.Tf(i18n.MsgErrorLoggerInit, err))
		os.Exit(1)
	}
	defer log.Close()

	// アプリケーション開始ログ
	log.Infof("%s", i18n.Tf(i18n.MsgAppStarted, i18n.Tf(i18n.MsgAppVersion, config.Version)))

	// 言語設定をパース
	cfg.Languages = ParseLanguages(languageStr)

	// バリデーション
	if cfg.AppID == "" && cfg.GameName == "" {
		fmt.Printf("%s\n\n", i18n.T(i18n.MsgErrorNoInput))
		printUsage()
		os.Exit(1)
	}

	if cfg.AppID != "" && cfg.GameName != "" {
		fmt.Printf("%s\n\n", i18n.T(i18n.MsgErrorBothInputs))
		printUsage()
		os.Exit(1)
	}

	// 出力ディレクトリの作成
	if cfg.OutputDir != "" {
		if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			log.Fatalf("%s", i18n.Tf(i18n.MsgErrorDirCreation, err))
		}
	}

	var reviews []models.ReviewData
	var appID string
	var gameName string
	var gameDetails *models.GameDetails

	// レビュー取得
	if cfg.AppID != "" {
		appID = cfg.AppID
		gameName = fmt.Sprintf("App ID %s", appID)
		log.Verbosef("App ID %s のレビューを取得中...", appID)
		reviews, err = api.FetchAllReviews(appID, cfg.MaxReviews, cfg.Verbose, cfg.Languages, cfg.Filter, log)
	} else {
		gameName = cfg.GameName
		log.Verbosef("ゲーム '%s' のレビューを取得中...", gameName)
		reviews, appID, err = api.GetReviewsByGameName(gameName, cfg.MaxReviews, cfg.Verbose, cfg.Languages, cfg.Filter, log)
	}

	if err != nil {
		log.Fatalf("%s", i18n.Tf(i18n.MsgErrorReviewFetch, err))
	}

	if len(reviews) == 0 {
		log.Info(i18n.T(i18n.MsgStatsNoReviews))
		return
	}

	log.Infof("取得したレビュー数: %d件", len(reviews))

	// ゲーム詳細情報を取得
	gameDetails, err = api.GetGameDetails(appID, cfg.Verbose, log)
	if err != nil {
		log.Verbosef("%s", i18n.Tf(i18n.MsgErrorGameDetailsInit, err))
		// ゲーム詳細情報が取得できなくてもレビュー保存は続行
		gameDetails = nil
	}

	// ファイル保存
	ext := ".txt"
	if cfg.OutputJSON {
		ext = ".json"
	}
	baseFilename := fmt.Sprintf("steam_reviews_%s%s", appID, ext)

	var savedFiles []string
	if cfg.SplitByLang {
		files, err := storage.SaveReviewsByLanguageWithGameDetails(reviews, baseFilename, cfg.OutputDir, cfg.Verbose, cfg.OutputJSON, gameDetails)
		if err != nil {
			log.Errorf("%s", i18n.Tf(i18n.MsgErrorFileSave, err))
		}
		savedFiles = files
	} else {
		filename := baseFilename
		if cfg.OutputDir != "" {
			filename = cfg.OutputDir + "/" + filename
		}

		if savedFile, err := storage.SaveReviewsToFileWithGameDetails(reviews, filename, cfg.OutputJSON, gameDetails); err != nil {
			log.Errorf("%s", i18n.Tf(i18n.MsgErrorFileSave, err))
		} else {
			savedFiles = append(savedFiles, savedFile)
			log.Verbosef("%s", i18n.Tf(i18n.MsgVerboseReviewSaved, filename))
		}
	}

	// 保存したファイル一覧を表示（標準出力のみ）
	log.Printf("\n%s", i18n.T(i18n.MsgFileSavedFiles))
	for _, file := range savedFiles {
		log.Printf("- %s\n", file)
	}
	log.Println()

	// ゲーム情報を使用して統計情報を表示
	displayGameName := gameName
	if gameDetails != nil {
		displayGameName = gameDetails.Name
	}

	// 統計情報を表示
	stats.PrintReviewStats(reviews, displayGameName, log)

	log.Info(i18n.T(i18n.MsgSuccessCompleted))
}
