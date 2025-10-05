package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type (
	FootballDataClient struct {
		competitionsMap map[Competition]string
		httpClient      *HttpClient
	}

	FootballDataResponse struct {
		Matches []struct {
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

func NewFootballDataClient() *FootballDataClient {
	return &FootballDataClient{
		httpClient: NewHttpClient(
			"https://api.football-data.org",
			WithHeader("x-auth-token", os.Getenv("AUTH_TOKEN")),
		),
		competitionsMap: map[Competition]string{
			// UefaChampionsLeague: "CL",
			// Bundesliga:          "BL1",
			// LaLiga:              "PD",
			// PremierLeague:       "PL",
			// Italiano:            "SA",
			// Brasileirao:         "BSA",
			Libertadores: "CLI",
		},
	}
}

func (c FootballDataClient) ListMatches(competition Competition, date time.Time) ([]Match, error) {
	competitionCode, found := c.competitionsMap[competition]
	if !found {
		return []Match{}, nil
	}

	currentDate := date.Format("2006-01-02")
	url := fmt.Sprintf("/v4/competitions/%s/matches?dateFrom=%s&dateTo=%s", competitionCode, currentDate, currentDate)
	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("list matches to %s: %v", competition, err)
	}

	var matchesResponse FootballDataResponse
	if err = json.Unmarshal(response, &matchesResponse); err != nil {
		return nil, fmt.Errorf("unmarshal matches to %s: %v", competition, err)
	}

	matches := make([]Match, len(matchesResponse.Matches))
	for i, match := range matchesResponse.Matches {
		matches[i] = Match{
			Competition: competition,
			StartAt:     match.UtcDate,
			HomeTeam:    match.HomeTeam.Name,
			AwayTeam:    match.AwayTeam.Name,
			HomeScore:   match.Score.FullTime.Home,
			AwayScore:   match.Score.FullTime.Away,
		}
	}
	return matches, nil
}

func (c *FootballDataClient) Contains(competition Competition) bool {
	_, found := c.competitionsMap[competition]
	return found
}

func (c *FootballDataClient) Name() string {
	return "FootballData"
}
