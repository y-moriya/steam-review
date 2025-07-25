package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/y-moriya/steam-review/internal/api"
	"github.com/y-moriya/steam-review/internal/models"
	"github.com/y-moriya/steam-review/internal/stats"
	"github.com/y-moriya/steam-review/internal/storage"
	"github.com/y-moriya/steam-review/pkg/config"
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
	name := config.AppName
	fmt.Printf(`%s version %s

使用方法:
  steam-review [オプション]

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
  steam-review -appid 440 -max 500 -verbose

  # 作成日時順でレビューを取得
  steam-review -appid 440 -max 500 -filter recent -verbose

  # ゲーム名で英語レビューを取得
  steam-review -game "Cyberpunk 2077" -lang "english" -max 1000 -output ./reviews

  # 複数言語のレビューを取得
  steam-review -game "Elden Ring" -lang "japanese,english" -max 300 -split

  # 日本語レビューをJSON形式で保存
  steam-review -appid 570 -max 2000 -output ./dota2_reviews -json -verbose

  # すべての言語のレビューを取得
  steam-review -appid 730 -lang "all" -max 1000 -split

  # 最近更新されたレビューから取得
  steam-review -appid 730 -filter updated -max 200

注意:
  - App IDとゲーム名のどちらか一方を指定してください
  - -lang を指定しない場合、デフォルトで日本語レビューのみを取得します
  - "all" を指定するとすべての言語のレビューを取得します
  - 大量のレビューを取得する場合は時間がかかります
  - Steam APIのレート制限により、リクエスト間に1秒の待機時間があります
`, name, config.Version)
}

func main() {
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
		fmt.Printf("%s version %s\n", config.AppName, config.Version)
		return
	}

	// ヘルプ表示
	if help {
		printUsage()
		return
	}

	// 言語設定をパース
	cfg.Languages = ParseLanguages(languageStr)

	// バリデーション
	if cfg.AppID == "" && cfg.GameName == "" {
		fmt.Printf("エラー: App ID またはゲーム名を指定してください\n\n")
		printUsage()
		os.Exit(1)
	}

	if cfg.AppID != "" && cfg.GameName != "" {
		fmt.Printf("エラー: App ID とゲーム名の両方を指定することはできません\n\n")
		printUsage()
		os.Exit(1)
	}

	// 出力ディレクトリの作成
	if cfg.OutputDir != "" {
		if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			log.Fatalf("出力ディレクトリの作成に失敗しました: %v", err)
		}
	}

	var reviews []models.ReviewData
	var appID string
	var gameName string
	var gameDetails *models.GameDetails
	var err error

	// レビュー取得
	if cfg.AppID != "" {
		appID = cfg.AppID
		gameName = fmt.Sprintf("App ID %s", appID)
		if cfg.Verbose {
			log.Printf("App ID %s のレビューを取得中...", appID)
		}
		reviews, err = api.FetchAllReviews(appID, cfg.MaxReviews, cfg.Verbose, cfg.Languages, cfg.Filter)
	} else {
		gameName = cfg.GameName
		if cfg.Verbose {
			log.Printf("ゲーム '%s' のレビューを取得中...", gameName)
		}
		reviews, appID, err = api.GetReviewsByGameName(gameName, cfg.MaxReviews, cfg.Verbose, cfg.Languages, cfg.Filter)
	}

	if err != nil {
		log.Fatalf("レビュー取得エラー: %v", err)
	}

	if len(reviews) == 0 {
		log.Printf("レビューが見つかりませんでした")
		return
	}

	// ゲーム詳細情報を取得
	gameDetails, err = api.GetGameDetails(appID, cfg.Verbose)
	if err != nil {
		if cfg.Verbose {
			log.Printf("ゲーム詳細情報の取得に失敗しました: %v", err)
		}
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
			log.Printf("ファイル保存エラー: %v", err)
		}
		savedFiles = files
	} else {
		filename := baseFilename
		if cfg.OutputDir != "" {
			filename = cfg.OutputDir + "/" + filename
		}

		if savedFile, err := storage.SaveReviewsToFileWithGameDetails(reviews, filename, cfg.OutputJSON, gameDetails); err != nil {
			log.Printf("ファイル保存エラー: %v", err)
		} else {
			savedFiles = append(savedFiles, savedFile)
			if cfg.Verbose {
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

	// ゲーム情報を使用して統計情報を表示
	displayGameName := gameName
	if gameDetails != nil {
		displayGameName = gameDetails.Name
	}

	// 統計情報を表示
	stats.PrintReviewStats(reviews, displayGameName)
}
