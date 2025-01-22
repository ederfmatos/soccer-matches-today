package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var httpClient = NewHttpClient("https://api.football-data.org", os.Getenv("AUTH_TOKEN"))

func main() {
	log.Println("Starting application")
	if err := run(); err != nil {
		_ = SendMessageToDiscord(fmt.Sprintf("Erro ao buscar os jogos de hoje: %s", err.Error()))
		log.Fatal(err)
	}
	log.Println("Finish application")
}

func run() error {

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
	matchesByCompetition := make(map[string][]Match)

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	channel := make(chan struct{}, 3)
	var errs []error

	for competitionCode, competitionName := range competitions {
		wg.Add(1)

		go func(code, name string) {
			channel <- struct{}{}
			defer func() {
				<-channel
				wg.Done()
			}()

			matches, err := listMatches(code, name)

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				errs = append(errs, err)
				return
			}
			matchesByCompetition[name] = matches
		}(competitionCode, competitionName)
	}
	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	message := createMessage(matchesByCompetition)
	if err := SendMessageToDiscord(message); err != nil {
		return fmt.Errorf("send message to discord: %v", err)
	}

	return nil
}

func listMatches(competitionCode, competitionName string) ([]Match, error) {
	log.Printf("Fetching matches to %s", competitionName)

	currentDate := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("/v4/competitions/%s/matches?dateFrom=%s&dateTo=%s", competitionCode, currentDate, currentDate)
	response, err := httpClient.Do(url)
	if err != nil {
		return nil, fmt.Errorf("list matches to %s: %v", competitionName, err)
	}

	var matchesResponse ListMatchesResponse
	if err = json.Unmarshal(response, &matchesResponse); err != nil {
		return nil, fmt.Errorf("unmarshal matches to %s: %v", competitionName, err)
	}
	defer log.Printf("Matches to %s fetched successfully", competitionName)

	matches := make([]Match, len(matchesResponse.Matches))
	for i, match := range matchesResponse.Matches {
		matches[i] = Match{
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
	return matches, nil
}

func createMessage(matchesByCompetition map[string][]Match) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# Jogos de Hoje - %s\n\n", time.Now().Format("02/01/2006")))

	if len(matchesByCompetition) == 0 {
		log.Println("No matches for today")
		return "Não existem jogos para hoje."
	}

	for competition, matches := range matchesByCompetition {
		if len(matches) == 0 {
			continue
		}

		builder.WriteString(fmt.Sprintf("## %s \n\n", competition))
		for _, match := range matches {
			builder.WriteString(createMessageToMatch(match))
		}
	}
	return builder.String()
}

func createMessageToMatch(match Match) string {
	startAt := match.StartAt.Format("15:04")
	if match.HomeScore == nil || match.AwayScore == nil {
		return fmt.Sprintf("**%s** x **%s** - :clock1: %s\n\n", match.HomeTeam, match.AwayTeam, startAt)
	}
	return fmt.Sprintf("**%s** %d x %d **%s** - :clock1: %s\n\n", match.HomeTeam, *match.HomeScore, *match.AwayScore, match.AwayTeam, startAt)
}
