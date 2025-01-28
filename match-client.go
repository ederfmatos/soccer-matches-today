package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
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

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	errs := make([]error, 0)

	for _, competition := range competitions {
		wg.Add(1)

		go func(competition Competition, wg *sync.WaitGroup) {
			defer wg.Done()

			for _, client := range m.clients {
				if !client.Contains(competition) {
					continue
				}

				log.Printf("Fetching Matches to %s using %v\n", competition, client.Name())

				matches, err := client.ListMatches(competition, currentDate)

				if err != nil {
					mutex.Lock()
					errs = append(errs, fmt.Errorf("list matches: %v", err))
					mutex.Unlock()
					continue
				}

				if len(matches) != 0 {
					mutex.Lock()
					matchesByCompetition[competition] = matches
					mutex.Unlock()
					return
				}
			}
		}(competition, &wg)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return matchesByCompetition, nil
}
