package i18n

import (
	"os"
	"strings"
)

// SupportedLanguages サポートしている言語のリスト
var SupportedLanguages = map[string]bool{
	"en": true,
	"ja": true,
}

// DetectLanguage 環境変数から言語を検出
func DetectLanguage() string {
	// 1. アプリケーション専用環境変数
	if lang := os.Getenv("STEAM_REVIEW_LANG"); lang != "" {
		if normalized := normalizeLanguage(lang); isSupportedLanguage(normalized) {
			return normalized
		}
	}

	// 2. 標準的な環境変数をチェック
	envVars := []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"}
	for _, envVar := range envVars {
		if lang := os.Getenv(envVar); lang != "" {
			if normalized := normalizeLanguage(lang); isSupportedLanguage(normalized) {
				return normalized
			}
		}
	}

	// 3. デフォルト言語
	return "en"
}

// normalizeLanguage 言語コードを正規化
func normalizeLanguage(lang string) string {
	// "ja_JP.UTF-8" -> "ja"
	lang = strings.ToLower(lang)
	parts := strings.FieldsFunc(lang, func(r rune) bool {
		return r == '_' || r == '.' || r == '-'
	})
	if len(parts) > 0 {
		return parts[0]
	}
	return lang
}

// isSupportedLanguage サポートされている言語かチェック
func isSupportedLanguage(lang string) bool {
	return SupportedLanguages[lang]
}
