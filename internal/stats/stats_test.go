package stats

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/y-moriya/steam-review/internal/models"
)

func TestPrintReviewStats(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

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
	PrintReviewStats(reviews, "Test Game")

	// 出力をキャプチャして復元
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// 出力を検証
	if !strings.Contains(output, "Test Game") {
		t.Error("Output should contain game name")
	}
	if !strings.Contains(output, "総レビュー数: 4") {
		t.Error("Output should contain total review count")
	}
	if !strings.Contains(output, "肯定的: 3 (75.0%)") {
		t.Error("Output should contain positive review stats")
	}
	if !strings.Contains(output, "否定的: 1 (25.0%)") {
		t.Error("Output should contain negative review stats")
	}
	if !strings.Contains(output, "japanese: 2件") {
		t.Error("Output should contain Japanese review count")
	}
	if !strings.Contains(output, "english: 2件") {
		t.Error("Output should contain English review count")
	}
}

func TestPrintReviewStats_EmptyReviews(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 空のレビューデータ
	reviews := []models.ReviewData{}

	// 関数を実行
	PrintReviewStats(reviews, "Empty Game")

	// 出力をキャプチャして復元
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// 出力を検証
	if !strings.Contains(output, "レビューが見つかりませんでした") {
		t.Error("Output should indicate no reviews found")
	}
}
