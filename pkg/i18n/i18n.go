package i18n

import "fmt"

// Localizer 国際化機能を提供
type Localizer struct {
	language string
	messages map[string]string
}

// Global localizer instance
var globalLocalizer *Localizer

// Init 国際化システムを初期化
func Init() {
	lang := DetectLanguage()
	globalLocalizer = &Localizer{
		language: lang,
		messages: getMessages(lang),
	}
}

// GetCurrentLanguage 現在の言語コードを取得
func GetCurrentLanguage() string {
	if globalLocalizer == nil {
		return "en"
	}
	return globalLocalizer.language
}

// T メッセージを翻訳（引数なし）
func T(key string) string {
	if globalLocalizer == nil {
		Init()
	}
	return globalLocalizer.T(key)
}

// Tf メッセージを翻訳（フォーマット引数あり）
func Tf(key string, args ...interface{}) string {
	if globalLocalizer == nil {
		Init()
	}
	return globalLocalizer.Tf(key, args...)
}

// T メッセージを翻訳（引数なし）
func (l *Localizer) T(key string) string {
	if msg, exists := l.messages[key]; exists {
		return msg
	}
	// フォールバック: 英語メッセージを試行
	if l.language != "en" {
		if enMsg, exists := getMessages("en")[key]; exists {
			return enMsg
		}
	}
	// 最終フォールバック: キー自体を返す
	return key
}

// Tf メッセージを翻訳（フォーマット引数あり）
func (l *Localizer) Tf(key string, args ...interface{}) string {
	template := l.T(key)
	return fmt.Sprintf(template, args...)
}

// getMessages 指定された言語のメッセージマップを取得
func getMessages(lang string) map[string]string {
	switch lang {
	case "ja":
		return getJapaneseMessages()
	case "en":
		return getEnglishMessages()
	default:
		return getEnglishMessages()
	}
}
