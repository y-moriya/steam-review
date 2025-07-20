package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCommandExamples(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Example 1: AppID with Japanese reviews",
			args:    []string{"-appid", "440", "-max", "500", "-verbose"},
			wantErr: false,
		},
		{
			name:    "Example 2: Game name with English reviews",
			args:    []string{"-game", "Cyberpunk 2077", "-lang", "english", "-max", "1000", "-output", filepath.Join(tempDir, "reviews")},
			wantErr: false,
		},
		{
			name:    "Example 3: Multiple languages",
			args:    []string{"-game", "Elden Ring", "-lang", "japanese,english", "-max", "300", "-split"},
			wantErr: false,
		},
		{
			name:    "Example 4: Japanese reviews in JSON format",
			args:    []string{"-appid", "570", "-max", "2000", "-output", filepath.Join(tempDir, "dota2_reviews"), "-json", "-verbose"},
			wantErr: false,
		},
		{
			name:    "Example 5: All languages",
			args:    []string{"-appid", "730", "-lang", "all", "-max", "1000", "-split"},
			wantErr: false,
		},
		{
			name:    "Example 6: Recently updated reviews",
			args:    []string{"-appid", "730", "-filter", "updated", "-max", "200"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// コマンドライン引数を一時的に設定
			oldArgs := os.Args
			os.Args = append([]string{"steam-review"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			// フラグをリセット
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			if err := runMain(); err != nil {
				if !tt.wantErr {
					t.Errorf("runMain() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

// テスト用にmain()をラップした関数
func runMain() error {
	// コマンドライン引数の定義
	var config Config
	var languageStr string
	var help bool
	var showVersion bool

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

	if showVersion {
		return nil
	}

	if help {
		return nil
	}

	config.Languages = parseLanguages(languageStr)

	if config.AppID == "" && config.GameName == "" {
		return fmt.Errorf("App ID またはゲーム名を指定してください")
	}

	if config.AppID != "" && config.GameName != "" {
		return fmt.Errorf("App ID とゲーム名の両方を指定することはできません")
	}

	if config.OutputDir != "" {
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			return fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
		}
	}

	var reviews []ReviewData
	var appID string
	var gameName string
	var err error

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
		return fmt.Errorf("レビュー取得エラー: %v", err)
	}

	if len(reviews) == 0 {
		return nil
	}

	ext := ".txt"
	if config.OutputJSON {
		ext = ".json"
	}
	baseFilename := fmt.Sprintf("steam_reviews_%s%s", appID, ext)

	if config.SplitByLang {
		_, err := SaveReviewsByLanguage(reviews, baseFilename, config.OutputDir, config.Verbose, config.OutputJSON)
		if err != nil {
			return fmt.Errorf("ファイル保存エラー: %v", err)
		}
	} else {
		filename := baseFilename
		if config.OutputDir != "" {
			filename = config.OutputDir + "/" + filename
		}

		_, err := SaveReviewsToFile(reviews, filename, config.OutputJSON)
		if err != nil {
			return fmt.Errorf("ファイル保存エラー: %v", err)
		}
	}

	return nil
}
