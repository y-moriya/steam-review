package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/y-moriya/steam-review/internal/api"
	"github.com/y-moriya/steam-review/internal/models"
	"github.com/y-moriya/steam-review/internal/storage"
	"github.com/y-moriya/steam-review/pkg/config"
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
			args:    []string{"-appid", "440", "-max", "5", "-verbose", "-output", filepath.Join(tempDir, "test1")},
			wantErr: false,
		},
		{
			name:    "Example 2: Game name with English reviews",
			args:    []string{"-game", "Cyberpunk 2077", "-lang", "english", "-max", "10", "-verbose", "-output", filepath.Join(tempDir, "test2")},
			wantErr: false,
		},
		{
			name:    "Example 3: Multiple languages",
			args:    []string{"-game", "Elden Ring", "-lang", "japanese,english", "-max", "3", "-verbose", "-split", "-output", filepath.Join(tempDir, "test3")},
			wantErr: false,
		},
		{
			name:    "Example 4: Japanese reviews in JSON format",
			args:    []string{"-appid", "570", "-max", "20", "-output", filepath.Join(tempDir, "test4"), "-json", "-verbose"},
			wantErr: false,
		},
		{
			name:    "Example 5: All languages",
			args:    []string{"-appid", "730", "-lang", "all", "-max", "200", "-verbose", "-split", "-output", filepath.Join(tempDir, "test5")},
			wantErr: false,
		},
		{
			name:    "Example 6: Recently updated reviews",
			args:    []string{"-appid", "730", "-filter", "updated", "-max", "2", "-verbose", "-output", filepath.Join(tempDir, "test6")},
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

func TestParseLanguages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Single language",
			input:    "japanese",
			expected: []string{"japanese"},
		},
		{
			name:     "Multiple languages",
			input:    "japanese,english,spanish",
			expected: []string{"japanese", "english", "spanish"},
		},
		{
			name:     "Languages with spaces",
			input:    "japanese, english, spanish",
			expected: []string{"japanese", "english", "spanish"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Only commas",
			input:    ",,,",
			expected: []string{},
		},
		{
			name:     "Mixed empty and valid",
			input:    "japanese,,english,",
			expected: []string{"japanese", "english"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLanguages(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseLanguages(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i, lang := range result {
				if lang != tt.expected[i] {
					t.Errorf("ParseLanguages(%q)[%d] = %q, want %q", tt.input, i, lang, tt.expected[i])
				}
			}
		})
	}
}

// テスト用にmain()をラップした関数
func runMain() error {
	// コマンドライン引数の定義
	var cfg config.Config
	var languageStr string
	var help bool
	var showVersion bool

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

	if showVersion {
		return nil
	}

	if help {
		return nil
	}

	cfg.Languages = ParseLanguages(languageStr)

	if cfg.AppID == "" && cfg.GameName == "" {
		return fmt.Errorf("App ID またはゲーム名を指定してください")
	}

	if cfg.AppID != "" && cfg.GameName != "" {
		return fmt.Errorf("App ID とゲーム名の両方を指定することはできません")
	}

	if cfg.OutputDir != "" {
		if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			return fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
		}
	}

	var reviews []models.ReviewData
	var appID string
	var gameName string
	var err error

	if cfg.AppID != "" {
		appID = cfg.AppID
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
		return fmt.Errorf("レビュー取得エラー: %v", err)
	}

	if len(reviews) == 0 {
		return nil
	}

	ext := ".txt"
	if cfg.OutputJSON {
		ext = ".json"
	}
	baseFilename := fmt.Sprintf("steam_reviews_%s%s", appID, ext)

	if cfg.SplitByLang {
		_, err := storage.SaveReviewsByLanguage(reviews, baseFilename, cfg.OutputDir, cfg.Verbose, cfg.OutputJSON)
		if err != nil {
			return fmt.Errorf("ファイル保存エラー: %v", err)
		}
	} else {
		filename := baseFilename
		if cfg.OutputDir != "" {
			filename = cfg.OutputDir + "/" + filename
		}

		_, err := storage.SaveReviewsToFile(reviews, filename, cfg.OutputJSON)
		if err != nil {
			return fmt.Errorf("ファイル保存エラー: %v", err)
		}
	}

	return nil
}
