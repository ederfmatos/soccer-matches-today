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

	Competition struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	ListMatchesResponse struct {
		Matches []struct {
			Competition struct {
				Name   string `json:"name"`
				Emblem string `json:"emblem"`
			} `json:"competition"`
			UtcDate  time.Time `json:"utcDate"`
			HomeTeam struct {
				Name string `json:"name"`
			} `json:"homeTeam"`
			AwayTeam struct {
				Name string `json:"name"`
			} `json:"awayTeam"`

			Score struct {
				FullTime struct {
					Home *int `json:"home"`
					Away *int `json:"away"`
				} `json:"fullTime"`
			} `json:"score"`
		} `json:"matches"`
	}
)
