package models

import (
	"encoding/json"
	"testing"
)

func TestFlexibleFloat64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "Number value",
			input:    `3.14`,
			expected: 3.14,
		},
		{
			name:     "String number",
			input:    `"2.5"`,
			expected: 2.5,
		},
		{
			name:     "Empty string",
			input:    `""`,
			expected: 0,
		},
		{
			name:     "Invalid string",
			input:    `"invalid"`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f FlexibleFloat64
			err := json.Unmarshal([]byte(tt.input), &f)
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}
			if float64(f) != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, want %v", float64(f), tt.expected)
			}
		})
	}
}

func TestConvertSteamReview(t *testing.T) {
	steamReview := SteamReview{
		RecommendationID:  "123456",
		Language:          "japanese",
		Review:            "Great game!",
		TimestampCreated:  1640995200,
		VotedUp:           true,
		VotesUp:           10,
		WeightedVoteScore: FlexibleFloat64(3.5),
		Author: SteamAuthor{
			SteamID:          "76561198000000000",
			NumGamesOwned:    100,
			PlayTimeForever:  120,
			PlayTimeAtReview: 60,
		},
	}

	result := ConvertSteamReview(steamReview)

	if result.RecommendationID != "123456" {
		t.Errorf("Expected RecommendationID to be '123456', got '%s'", result.RecommendationID)
	}
	if result.Language != "japanese" {
		t.Errorf("Expected Language to be 'japanese', got '%s'", result.Language)
	}
	if result.Review != "Great game!" {
		t.Errorf("Expected Review to be 'Great game!', got '%s'", result.Review)
	}
	if !result.VotedUp {
		t.Errorf("Expected VotedUp to be true, got %v", result.VotedUp)
	}
	if result.WeightedScore != 3.5 {
		t.Errorf("Expected WeightedScore to be 3.5, got %v", result.WeightedScore)
	}
	if result.Author.SteamID != "76561198000000000" {
		t.Errorf("Expected Author.SteamID to be '76561198000000000', got '%s'", result.Author.SteamID)
	}
}
