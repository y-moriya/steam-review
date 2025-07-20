func main() {
	// 使用例1: App IDを直接指定
	appID := "2089600"
	maxReviews := 5_000_000
	
	log.Printf("App ID %s のレビューを取得中...", appID)
	reviews, err := FetchAllReviews(appID, maxReviews)
	if err != nil {
		log.Printf("レビュー取得エラー: %v", err)
	} else {
		// 言語別にファイル保存
		baseFilename := fmt.Sprintf("steam_reviews_%s.json", appID)
		if err := SaveReviewsByLanguage(reviews, baseFilename); err != nil {
			log.Printf("ファイル保存エラー: %v", err)
		}
		PrintReviewStats(reviews, "都市伝説解体センター")
	}

}