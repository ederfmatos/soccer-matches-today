package main

import (
	"fmt"
	"github.com/ederfmatos/go-concurrency/pkg/concurrency"
	"log"
	"time"
)

type (
	MatchClientDelegator struct {
		clients []MatchesClient
	}

	MatchesClient interface {
		ListMatches(competition Competition, date time.Time) ([]Match, error)
		Contains(competition Competition) bool
		Name() string
	}
)

func NewMatchClient() *MatchClientDelegator {
	return &MatchClientDelegator{
		clients: []MatchesClient{
			NewUOLClient(),
			NewFootballDataClient(),
		},
	}
}

func (m *MatchClientDelegator) ListMatches() (map[Competition][]Match, error) {
	matchesByCompetition := make(map[Competition][]Match)
	currentDate := time.Now().In(brazilLocation)

	matches, err := concurrency.ForEach[Competition, []Match](competitions, 4, func(competition Competition) ([]Match, error) {
		for _, client := range m.clients {
			if !client.Contains(competition) {
				continue
			}
			log.Printf("Fetching Matches to %s using %v\n", competition, client.Name())

			matches, err := client.ListMatches(competition, currentDate)

			if err != nil {
				return nil, fmt.Errorf("list matches: %v", err)
			}
			return matches, nil
		}

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	for _, items := range matches {
		if len(items) > 0 {
			matchesByCompetition[items[0].Competition] = items
		}
	}

	return matchesByCompetition, nil
}
