# Steam Reviews CLI Tool

[![Test](https://github.com/y-moriya/steam-review/actions/workflows/test.yml/badge.svg)](https://github.com/y-moriya/steam-review/actions/workflows/test.yml)

English | [日本語](README.ja.md)

Steam Reviews is a command-line tool for retrieving and saving Steam game reviews.
You can specify a game by its App ID or name and save the reviews in either JSON or text format.

## Features

- Search games by App ID or game name
- Filter reviews by language
- Specify maximum number of reviews to retrieve
- Sort reviews by creation date or update time
- Save in JSON or text format
- Split files by language
- Display detailed statistics

## Installation

```bash
go install github.com/y-moriya/steam-review@latest
```

## Usage

```
steam-review [options]
```

### Options

| Option     | Description | Default |
|------------|-------------|---------|
| -appid     | Steam App ID (e.g., 440) | - |
| -game      | Game name (e.g., "Team Fortress 2") | - |
| -max       | Maximum number of reviews to retrieve (0 for unlimited) | 100 |
| -lang      | Language to retrieve (comma-separated) | japanese |
| -output    | Output directory | output |
| -split     | Split files by language | false |
| -json      | Save output files in JSON format (.json) | false |
| -filter    | Review filter (recent/updated/all) | all |
| -verbose   | Display detailed logs | false |
| -help      | Display help | false |
| -version   | Display version information | - |

### Filter Options

- `all`: Sort by helpfulness (default)
- `recent`: Sort by creation date
- `updated`: Sort by last update time

## Examples

1. Get Japanese reviews by App ID (default)
```bash
steam-review -appid 440 -max 500 -verbose
```

2. Get English reviews by game name
```bash
steam-review -game "Cyberpunk 2077" -lang "english" -max 1000 -output ./reviews
```

3. Get reviews in multiple languages
```bash
steam-review -game "Elden Ring" -lang "japanese,english" -max 300 -split
```

4. Save Japanese reviews in JSON format
```bash
steam-review -appid 570 -max 2000 -output ./dota2_reviews -json -verbose
```

5. Get reviews in all languages
```bash
steam-review -appid 730 -lang "all" -max 1000 -split
```

6. Get recently updated reviews
```bash
steam-review -appid 730 -filter updated -max 200
```

## Output Files

### Text Format (Default)

```
=== Review 1 ===
ID: 195635539
language: japanese
voted_up: true
votes_up: 54
votes_funny: 28
weighted_score: 0.82
steam_purchase: true
playtime: 3754 minutes
created_at: 2025-05-26 00:49:30
updated_at: 2025-05-26 22:12:04
review:
Review text
```

### JSON形式 (-json option)

```json
[
  {
    "recommendation_id": "12345678",
    "author": {
      "steam_id": "76561197960287930",
      "num_games_owned": 100,
      "num_reviews": 10,
      "playtime_forever": 1000,
      "playtime_last_two_weeks": 10,
      "playtime_at_review": 500,
      "last_played": 1626825600
    },
    "language": "japanese",
    "review": "レビュー本文",
    "timestamp_created": 1626825600,
    "timestamp_updated": 1626825600,
    "voted_up": true,
    "votes_up": 10,
    "votes_funny": 5,
    "weighted_vote_score": 0.8,
    "comment_count": 3,
    "steam_purchase": true,
    "received_for_free": false,
    "written_during_early_access": false
  }
]
```

## Notes

- Specify either App ID or game name, not both
- If `-lang` is not specified, only Japanese reviews will be retrieved by default
- Use `all` to retrieve reviews in all languages
- Retrieving a large number of reviews may take time
- Due to Steam API rate limits, there is a 1-second delay between requests
- Output directory will be created automatically

## License

[MIT License](LICENSE)
