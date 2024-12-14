package cmd

import (
	dm "github.com/ricardosilva86/pokedexcli/datamodels"
	cache "github.com/ricardosilva86/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	Name        string
	Description string
	Command     func(config *dm.Config, pokemonCache *cache.Cache, value string) error
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Display all command's help",
			Command:     CommandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "exit the application",
			Command:     CommandExit,
		},
		"map": {
			Name:        "map",
			Description: "Show all location areas in the map",
			Command:     CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Show all location areas in the map but backwards",
			Command:     CommandMapB,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a location area",
			Command:     CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a Pokemon",
			Command:     CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a Pokemon",
			Command:     CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Show your pokedex",
			Command:     CommandPokedex,
		},
	}
}
