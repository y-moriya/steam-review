package config

const (
	// バージョン情報
	Version = "v0.5.1"                 // プログラムのバージョン
	AppName = "Steam Reviews CLI Tool" // プログラム名

	// レビューのフィルター
	FilterAll     = "all"     // 有用性による並び替え
	FilterRecent  = "recent"  // 作成日時による並び替え
	FilterUpdated = "updated" // 最終更新日時による並び替え

	// ファイル形式
	FileExtJSON = ".json" // JSON形式のファイル拡張子
	FileExtTXT  = ".txt"  // テキスト形式のファイル拡張子
)

// Config コマンドライン引数の設定
type Config struct {
	AppID       string
	GameName    string
	MaxReviews  int
	Languages   []string
	OutputDir   string
	Verbose     bool
	SplitByLang bool
	OutputJSON  bool
	Filter      string // レビューのフィルター
}
