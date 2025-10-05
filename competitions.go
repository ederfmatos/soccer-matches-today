package main

type Competition string

const (
	Paulistao                  Competition = "Campeonato Paulista"
	UefaChampionsLeague        Competition = "UEFA Champions League"
	Saudita                    Competition = "Campeonato Saudita"
	Bundesliga                 Competition = "Bundesliga"
	LaLiga                     Competition = "La Liga"
	PremierLeague              Competition = "Premier League"
	Italiano                   Competition = "Serie A Italiana"
	Libertadores               Competition = "Libertadores"
	Brasileirao                Competition = "Brasileirão Serie A"
	EliminatoriasEuropeias     Competition = "Eliminatorias Europeias"
	AmistososSelecaoBrasileira Competition = "Amistosos Seleção Brasileira"
)

var competitions = []Competition{
	Paulistao,
	UefaChampionsLeague,
	Saudita,
	Bundesliga,
	LaLiga,
	PremierLeague,
	Italiano,
	Libertadores,
	Brasileirao,
	EliminatoriasEuropeias,
	AmistososSelecaoBrasileira,
}
