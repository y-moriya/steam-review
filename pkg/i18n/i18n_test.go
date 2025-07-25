package i18n

import (
	"os"
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "STEAM_REVIEW_LANG takes priority",
			envVars:  map[string]string{"STEAM_REVIEW_LANG": "ja", "LANG": "en"},
			expected: "ja",
		},
		{
			name:     "Fallback to LANG",
			envVars:  map[string]string{"LANG": "ja_JP.UTF-8"},
			expected: "ja",
		},
		{
			name:     "Fallback to LC_ALL",
			envVars:  map[string]string{"LC_ALL": "en_US.UTF-8"},
			expected: "en",
		},
		{
			name:     "Default to English",
			envVars:  map[string]string{},
			expected: "en",
		},
		{
			name:     "Unsupported language defaults to English",
			envVars:  map[string]string{"STEAM_REVIEW_LANG": "fr"},
			expected: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 環境変数をクリア
			clearEnvVars()

			// テスト用環境変数を設定
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// テスト実行
			result := DetectLanguage()

			// 検証
			if result != tt.expected {
				t.Errorf("DetectLanguage() = %v, want %v", result, tt.expected)
			}

			// 環境変数をクリア
			clearEnvVars()
		})
	}
}

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ja", "ja"},
		{"ja_JP", "ja"},
		{"ja_JP.UTF-8", "ja"},
		{"en-US", "en"},
		{"EN_US", "en"},
		{"C", "c"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeLanguage(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeLanguage(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSupportedLanguage(t *testing.T) {
	tests := []struct {
		lang     string
		expected bool
	}{
		{"en", true},
		{"ja", true},
		{"fr", false},
		{"", false},
		{"zh", false},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			result := isSupportedLanguage(tt.lang)
			if result != tt.expected {
				t.Errorf("isSupportedLanguage(%q) = %v, want %v", tt.lang, result, tt.expected)
			}
		})
	}
}

func TestLocalizerT(t *testing.T) {
	tests := []struct {
		language string
		key      string
		expected string
	}{
		{"en", "error.no_input", "Error: Please specify either App ID or game name"},
		{"ja", "error.no_input", "エラー: App ID またはゲーム名を指定してください"},
		{"en", "stats.title", "=== Review Statistics ==="},
		{"ja", "stats.title", "=== レビュー統計 ==="},
		{"en", "nonexistent.key", "nonexistent.key"}, // フォールバック
	}

	for _, tt := range tests {
		t.Run(tt.language+"_"+tt.key, func(t *testing.T) {
			localizer := &Localizer{
				language: tt.language,
				messages: getMessages(tt.language),
			}

			result := localizer.T(tt.key)
			if result != tt.expected {
				t.Errorf("Localizer.T(%q) = %q, want %q", tt.key, result, tt.expected)
			}
		})
	}
}

func TestLocalizerTf(t *testing.T) {
	tests := []struct {
		language string
		key      string
		args     []interface{}
		expected string
	}{
		{"en", "stats.total_reviews", []interface{}{100}, "Total reviews: 100"},
		{"ja", "stats.total_reviews", []interface{}{100}, "総レビュー数: 100"},
		{"en", "stats.positive", []interface{}{80, 80.0}, "Positive: 80 (80.0%)"},
		{"ja", "stats.positive", []interface{}{80, 80.0}, "肯定的: 80 (80.0%)"},
	}

	for _, tt := range tests {
		t.Run(tt.language+"_"+tt.key, func(t *testing.T) {
			localizer := &Localizer{
				language: tt.language,
				messages: getMessages(tt.language),
			}

			result := localizer.Tf(tt.key, tt.args...)
			if result != tt.expected {
				t.Errorf("Localizer.Tf(%q, %v) = %q, want %q", tt.key, tt.args, result, tt.expected)
			}
		})
	}
}

func TestGlobalFunctions(t *testing.T) {
	// テスト前に環境変数をクリア
	clearEnvVars()

	// 英語環境でテスト
	os.Setenv("STEAM_REVIEW_LANG", "en")
	globalLocalizer = nil // リセット

	Init()

	if GetCurrentLanguage() != "en" {
		t.Errorf("GetCurrentLanguage() = %v, want 'en'", GetCurrentLanguage())
	}

	result := T("error.no_input")
	expected := "Error: Please specify either App ID or game name"
	if result != expected {
		t.Errorf("T('error.no_input') = %q, want %q", result, expected)
	}

	result = Tf("stats.total_reviews", 50)
	expected = "Total reviews: 50"
	if result != expected {
		t.Errorf("Tf('stats.total_reviews', 50) = %q, want %q", result, expected)
	}

	// 日本語環境でテスト
	os.Setenv("STEAM_REVIEW_LANG", "ja")
	globalLocalizer = nil // リセット

	Init()

	if GetCurrentLanguage() != "ja" {
		t.Errorf("GetCurrentLanguage() = %v, want 'ja'", GetCurrentLanguage())
	}

	result = T("error.no_input")
	expected = "エラー: App ID またはゲーム名を指定してください"
	if result != expected {
		t.Errorf("T('error.no_input') = %q, want %q", result, expected)
	}

	// クリーンアップ
	clearEnvVars()
	globalLocalizer = nil
}

func TestFallbackBehavior(t *testing.T) {
	// 日本語環境で存在しないキーをテスト
	os.Setenv("STEAM_REVIEW_LANG", "ja")
	globalLocalizer = nil

	Init()

	// 存在しないキーの場合、英語にフォールバックするかキー自体を返す
	result := T("nonexistent.key")
	if result != "nonexistent.key" {
		t.Errorf("T('nonexistent.key') = %q, want 'nonexistent.key'", result)
	}

	clearEnvVars()
	globalLocalizer = nil
}

// clearEnvVars テスト用の環境変数をクリア
func clearEnvVars() {
	envVars := []string{"STEAM_REVIEW_LANG", "LANG", "LC_ALL", "LC_MESSAGES", "LANGUAGE"}
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
