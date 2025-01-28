package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

func main() {
	log.Println("Starting application")
	notificator := NewNotificator()
	if err := run(notificator); err != nil {
		_ = notificator.SendMessage(fmt.Sprintf("Erro ao buscar os jogos de hoje: %s", err.Error()))
		log.Fatal(err)
	}
	log.Println("Finish application")
}

func run(notificator Notificator) error {
	matchesClient := NewMatchClient()
	matchesByCompetition, err := matchesClient.ListMatches()
	if err != nil {
		return fmt.Errorf("list matches: %v", err)
	}

	message := createMessage(matchesByCompetition)
	if err = notificator.SendMessage(message); err != nil {
		return fmt.Errorf("send message: %v", err)
	}

	return nil
}

func createMessage(matchesByCompetition map[Competition][]Match) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# Jogos de Hoje - %s \n\n", time.Now().Format("02/01/2006")))

	matchesCount := 0
	for competition, matches := range matchesByCompetition {
		if len(matches) == 0 {
			continue
		}
		matchesCount += len(matches)

		sort.Slice(matches, func(i, j int) bool {
			return matches[i].StartAt.Before(matches[j].StartAt)
		})

		builder.WriteString(fmt.Sprintf("## %s \n\n", competition))
		for _, match := range matches {
			builder.WriteString(createMessageToMatch(match))
		}
	}

	if matchesCount == 0 {
		builder.WriteString("Não existem jogos para hoje")
	}

	return builder.String()
}

func createMessageToMatch(match Match) string {
	startAt := match.StartAt.Format("15h04")
	if match.HomeScore == nil || match.AwayScore == nil {
		return fmt.Sprintf("%s - **%s** x **%s**\n", startAt, match.HomeTeam, match.AwayTeam)
	}
	return fmt.Sprintf("%s - **%s** %d x %d **%s**\n", startAt, match.HomeTeam, *match.HomeScore, *match.AwayScore, match.AwayTeam)
}
