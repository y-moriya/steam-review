package config

import "testing"

func TestConstants(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if AppName == "" {
		t.Error("AppName should not be empty")
	}
	if FilterAll != "all" {
		t.Errorf("FilterAll should be 'all', got '%s'", FilterAll)
	}
	if FilterRecent != "recent" {
		t.Errorf("FilterRecent should be 'recent', got '%s'", FilterRecent)
	}
	if FilterUpdated != "updated" {
		t.Errorf("FilterUpdated should be 'updated', got '%s'", FilterUpdated)
	}
	if FileExtJSON != ".json" {
		t.Errorf("FileExtJSON should be '.json', got '%s'", FileExtJSON)
	}
	if FileExtTXT != ".txt" {
		t.Errorf("FileExtTXT should be '.txt', got '%s'", FileExtTXT)
	}
}

func TestConfigStruct(t *testing.T) {
	cfg := Config{
		AppID:       "440",
		GameName:    "Team Fortress 2",
		MaxReviews:  100,
		Languages:   []string{"japanese", "english"},
		OutputDir:   "./output",
		Verbose:     true,
		SplitByLang: false,
		OutputJSON:  true,
		Filter:      FilterAll,
	}

	if cfg.AppID != "440" {
		t.Errorf("Expected AppID to be '440', got '%s'", cfg.AppID)
	}
	if cfg.MaxReviews != 100 {
		t.Errorf("Expected MaxReviews to be 100, got %d", cfg.MaxReviews)
	}
	if len(cfg.Languages) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(cfg.Languages))
	}
	if !cfg.Verbose {
		t.Error("Expected Verbose to be true")
	}
	if !cfg.OutputJSON {
		t.Error("Expected OutputJSON to be true")
	}
}
