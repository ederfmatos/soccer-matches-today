package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	log.Println("Starting application")
	if err := run(); err != nil {
		_ = SendMessageToDiscord(fmt.Sprintf("Erro ao buscar os jogos de hoje: %s", err.Error()))
		log.Fatal(err)
	}
	log.Println("Finish application")
}

func run() error {
	httpClient := NewHttpClient("https://api.football-data.org", os.Getenv("AUTH_TOKEN"))

	currentDate := time.Now().Add(time.Hour * -24).Format("2006-01-02")
	matchesByCompetition := make(map[string][]Match)
	competitions := map[string]string{
		"CL":  "UEFA Champions League",
		"BL1": "Bundesliga",
		"BSA": "Brasileirão Serie A",
		"PD":  "La Liga",
		"FL1": "Ligue 1",
		"PPL": "Liga Portuguesa",
		"SA":  "Serie A Italiana",
		"PL":  "Premier League",
		"CLI": "Copa Libertadores",
	}
	for competitionCode, competitionName := range competitions {
		url := fmt.Sprintf("/v4/competitions/%s/matches?dateFrom=%s&dateTo=%s", competitionCode, currentDate, currentDate)
		response, err := httpClient.Do(url)
		if err != nil {
			return fmt.Errorf("list matches to %s: %v", competitionName, err)
		}

		var matchesResponse ListMatchesResponse
		if err = json.Unmarshal(response, &matchesResponse); err != nil {
			return fmt.Errorf("unmarshal matches to %s: %v", competitionName, err)
		}
		log.Printf("Matches to %s fetched successfully", competitionName)

		matchesByCompetition[competitionName] = make([]Match, len(matchesResponse.Matches))
		for i, match := range matchesResponse.Matches {
			matchesByCompetition[competitionName][i] = Match{
				Competition: Competition{
					Name:  match.Competition.Name,
					Image: match.Competition.Emblem,
				},
				StartAt:   match.UtcDate.In(time.Local),
				HomeTeam:  match.HomeTeam.Name,
				AwayTeam:  match.AwayTeam.Name,
				HomeScore: match.Score.FullTime.Home,
				AwayScore: match.Score.FullTime.Away,
			}
		}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# Jogos de Hoje - %s\n\n", time.Now().Add(time.Hour*-24).Format("02/01/2006")))

	if len(matchesByCompetition) == 0 {
		log.Println("No matches for today")

		builder.WriteString("Não existem jogos para hoje.")
		if err := SendMessageToDiscord(builder.String()); err != nil {
			return fmt.Errorf("send message to discord: %v", err)
		}
		return nil
	}

	for competition, matches := range matchesByCompetition {
		builder.WriteString("## ")
		builder.WriteString(competition)
		builder.WriteString("\n\n")

		for _, match := range matches {
			if match.HomeScore != nil && match.AwayScore != nil {
				builder.WriteString(fmt.Sprintf("%s %d x %d %s - %s\n\n", match.HomeTeam, *match.HomeScore, *match.AwayScore, match.AwayTeam, match.StartAt.Format("15:04")))
			} else {
				builder.WriteString(fmt.Sprintf("%s x %s - %s\n\n", match.HomeTeam, match.AwayTeam, match.StartAt.Format("15:04")))
			}
		}
	}

	log.Println("Sending matches to discord")
	if err := SendMessageToDiscord(builder.String()); err != nil {
		return fmt.Errorf("send message to discord: %v", err)
	}

	return nil
}
