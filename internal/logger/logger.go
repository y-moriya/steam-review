package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger カスタムロガー
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
	verbose     bool
}

// New 新しいロガーを作成
func New(logDir string, verbose bool) (*Logger, error) {
	// ログディレクトリを作成
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("ログディレクトリの作成に失敗しました: %v", err)
	}

	// ログファイル名を現在時刻で生成
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(logDir, fmt.Sprintf("steam-review_%s.log", timestamp))

	// ログファイルを開く
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("ログファイルのオープンに失敗しました: %v", err)
	}

	// マルチライター（標準出力とファイルの両方に出力）
	multiWriter := io.MultiWriter(os.Stdout, file)
	multiErrorWriter := io.MultiWriter(os.Stderr, file)

	// ロガーを作成
	infoLogger := log.New(multiWriter, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	errorLogger := log.New(multiErrorWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	return &Logger{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		file:        file,
		verbose:     verbose,
	}, nil
}

// Info 情報ログを出力
func (l *Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

// Infof フォーマット付き情報ログを出力
func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error エラーログを出力
func (l *Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

// Errorf フォーマット付きエラーログを出力
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Fatal 致命的エラーログを出力してプログラムを終了
func (l *Logger) Fatal(v ...interface{}) {
	l.errorLogger.Fatalln(v...)
}

// Fatalf フォーマット付き致命的エラーログを出力してプログラムを終了
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.errorLogger.Fatalf(format, v...)
}

// Verbose 詳細ログを出力（verboseフラグが有効な場合のみ）
func (l *Logger) Verbose(v ...interface{}) {
	if l.verbose {
		l.infoLogger.Println(v...)
	}
}

// Verbosef フォーマット付き詳細ログを出力（verboseフラグが有効な場合のみ）
func (l *Logger) Verbosef(format string, v ...interface{}) {
	if l.verbose {
		l.infoLogger.Printf(format, v...)
	}
}

// Print 標準出力のみに出力（ログファイルには出力しない）
func (l *Logger) Print(v ...interface{}) {
	fmt.Print(v...)
}

// Printf フォーマット付きで標準出力のみに出力（ログファイルには出力しない）
func (l *Logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Println 標準出力のみに出力（ログファイルには出力しない）
func (l *Logger) Println(v ...interface{}) {
	fmt.Println(v...)
}

// Close ログファイルをクローズ
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
