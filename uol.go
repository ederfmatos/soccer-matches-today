package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type (
	UOLClient struct {
		competitionsMap map[string]Competition
		httpClient      *HttpClient
		uolResponse     *UOLResponse
		mutex           sync.Mutex
	}

	UOLResponse struct {
		Teams map[string]struct {
			Name      string `json:"nome-completo"`
			ShortName string `json:"sigla"`
			Emblem    string `json:"brasao"`
		} `json:"equipes"`

		Games map[string]struct {
			Competition   string `json:"competicao"`
			CompetitionID string `json:"id-competicao"`
			Date          string `json:"data"`
			Hour          string `json:"horario"`
			HomeTeam      string `json:"time1"`
			AwayTeam      string `json:"time2"`
			HomeScore     *int   `json:"placar1"`
			AwayScore     *int   `json:"placar2"`
		} `json:"jogos"`
	}
)

func NewUOLClient() *UOLClient {
	return &UOLClient{
		httpClient: NewHttpClient("https://www.uol.com.br/esporte/service/"),
		competitionsMap: map[string]Competition{
			"104": Paulistao,
			"83":  UefaChampionsLeague,
			"178": Saudita,
			"12":  Bundesliga,
			"72":  LaLiga,
			"79":  PremierLeague,
			"81":  Italiano,
		},
	}
}

func (c *UOLClient) ListMatches(competition Competition, date time.Time) ([]Match, error) {
	c.mutex.Lock()

	if c.uolResponse == nil {
		response, err := c.httpClient.Get(`?loadComponent=api&data={"module":"tools","api":"json","method":"open","busca":"commons.uol.com.br/sistemas/esporte/modalidades/futebol/campeonatos/etc/jogos/resultados_e_proximos/dados.json"}`)
		if err != nil {
			return nil, fmt.Errorf("list matches: %v", err)
		}

		var uolResponse UOLResponse
		if err = json.Unmarshal(response, &uolResponse); err != nil {
			return nil, fmt.Errorf("unmarshal matches: %v", err)
		}

		c.uolResponse = &uolResponse
	}

	c.mutex.Unlock()

	matches := make([]Match, 0)
	for _, game := range c.uolResponse.Games {
		gameCompetition, found := c.competitionsMap[game.CompetitionID]
		if !found || gameCompetition != competition {
			continue
		}

		rawTime := fmt.Sprintf("%s %s", game.Date, game.Hour)
		startAt, err := time.Parse("2006-01-02 15h04", rawTime)
		if err != nil {
			return nil, fmt.Errorf("parse time: %s = %v", rawTime, err)
		}

		startYear, startMonth, startDay := startAt.Date()
		currentYear, currentMonth, currentDay := date.Date()
		if startYear != currentYear || startMonth != currentMonth || startDay != currentDay {
			continue
		}

		matches = append(matches, Match{
			Competition: competition,
			StartAt:     startAt.UTC(),
			HomeTeam:    c.uolResponse.Teams[game.HomeTeam].Name,
			AwayTeam:    c.uolResponse.Teams[game.AwayTeam].Name,
			HomeScore:   game.HomeScore,
			AwayScore:   game.AwayScore,
		})
	}

	return matches, nil
}

func (c *UOLClient) Contains(competition Competition) bool {
	for _, value := range c.competitionsMap {
		if value == competition {
			return true
		}
	}
	return false
}

func (c *UOLClient) Name() string {
	return "UOL"
}
