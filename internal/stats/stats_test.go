package stats

import (
	"strings"
	"testing"

	"github.com/y-moriya/steam-review/internal/models"
)

// TestLogger テスト用のロガー実装
type TestLogger struct {
	output []string
}

// NewTestLogger テスト用ロガーを作成
func NewTestLogger() *TestLogger {
	return &TestLogger{
		output: make([]string, 0),
	}
}

// Println ログを出力
func (l *TestLogger) Println(v ...interface{}) {
	var parts []string
	for _, val := range v {
		parts = append(parts, strings.TrimSpace(val.(string)))
	}
	l.output = append(l.output, strings.Join(parts, " "))
}

// Printf フォーマット付きログを出力
func (l *TestLogger) Printf(format string, v ...interface{}) {
	// fmt.Sprintf相当の処理を簡略化
	result := format
	for _, val := range v {
		switch val := val.(type) {
		case string:
			result = strings.Replace(result, "%s", val, 1)
		case int:
			result = strings.Replace(result, "%d", "4", 1) // テスト用固定値
		case float64:
			result = strings.Replace(result, "%.1f%%", "75.0%", 1) // テスト用固定値
		}
	}
	l.output = append(l.output, result)
}

// GetOutput 出力された内容を結合して取得
func (l *TestLogger) GetOutput() string {
	return strings.Join(l.output, "\n")
}

func TestPrintReviewStats(t *testing.T) {
	// テスト用ロガーを作成
	logger := NewTestLogger()

	// テストデータ
	reviews := []models.ReviewData{
		{
			Language: "japanese",
			VotedUp:  true,
		},
		{
			Language: "japanese",
			VotedUp:  false,
		},
		{
			Language: "english",
			VotedUp:  true,
		},
		{
			Language: "english",
			VotedUp:  true,
		},
	}

	// 関数を実行
	PrintReviewStats(reviews, "Test Game", logger)

	// 出力を取得
	output := logger.GetOutput()

	// 出力を検証
	if !strings.Contains(output, "Test Game") {
		t.Error("Output should contain game name")
	}
	if !strings.Contains(output, "総レビュー数") {
		t.Error("Output should contain total review count")
	}
	if !strings.Contains(output, "肯定的") {
		t.Error("Output should contain positive review stats")
	}
	if !strings.Contains(output, "否定的") {
		t.Error("Output should contain negative review stats")
	}
	if !strings.Contains(output, "言語別レビュー統計") {
		t.Error("Output should contain language statistics")
	}
}

func TestPrintReviewStats_EmptyReviews(t *testing.T) {
	// テスト用ロガーを作成
	logger := NewTestLogger()

	// 空のレビューデータ
	reviews := []models.ReviewData{}

	// 関数を実行
	PrintReviewStats(reviews, "Empty Game", logger)

	// 出力を取得
	output := logger.GetOutput()

	// 出力を検証
	if !strings.Contains(output, "レビューが見つかりませんでした") {
		t.Error("Output should indicate no reviews found")
	}
}
