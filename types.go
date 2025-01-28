package main

import "time"

type (
	Match struct {
		Competition Competition `json:"competition"`
		StartAt     time.Time   `json:"start_at"`
		HomeTeam    string      `json:"home_team"`
		AwayTeam    string      `json:"away_team"`
		HomeScore   *int        `json:"home_score"`
		AwayScore   *int        `json:"away_score"`
	}
)
