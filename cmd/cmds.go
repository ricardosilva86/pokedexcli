package cmd

import (
	"encoding/json"
	"fmt"
	dm "github.com/ricardosilva86/pokedexcli/datamodels"
	internal "github.com/ricardosilva86/pokedexcli/internal/config"
	cache "github.com/ricardosilva86/pokedexcli/internal/pokecache"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	podex = dm.NewPokeDex()
)

func apiCall(config *dm.Config, back bool) ([]byte, error) {
	if back {
		if config.PreviousURL == nil {
			config.URL = &internal.BaseURL
		}
		config.URL = config.PreviousURL
	}
	req, err := http.NewRequest("GET", *config.URL, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("error creating request: %v\n", err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error sending request: %v\n", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading response body: %v\n", err)
	}

	if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("error calling api: %v\n", res.Status)
	}

	return data, nil
}
func CommandExit(config *dm.Config, pokemonCache *cache.Cache, location string) error {
	os.Exit(0)
	return nil
}

func CommandHelp(config *dm.Config, pokemonCache *cache.Cache, location string) error {
	fmt.Println("How to use this program:")
	for _, cmd := range GetCommands() {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

func CommandMap(config *dm.Config, pokemonCache *cache.Cache, location string) error {
	if config.URL == nil || *config.URL == "" {
		return fmt.Errorf("no url found or url is empty\n")
	}

	data := []byte{}
	if d, ok := pokemonCache.Get(*config.URL); ok {
		fmt.Println("cache hit")
		data = d
	} else {
		var err error
		data, err = apiCall(config, false)
		if err != nil {
			return err
		}
		pokemonCache.Add(*config.URL, data)
	}
	var locations dm.Location
	if err := json.Unmarshal(data, &locations); err != nil {
		return fmt.Errorf("error unmarshalling response body: %v\n", err)
	}
	for _, localtion := range locations.Results {
		fmt.Println(localtion.Name)
	}

	//set next url
	config.URL = &locations.Next
	if locations.Next == "nil" {
		config.URL = &internal.BaseURL
	}
	//set previous url
	if locations.Previous != "null" {
		config.PreviousURL = &locations.Previous
	}

	return nil
}

func CommandMapB(config *dm.Config, pokemonCache *cache.Cache, location string) error {
	if config.URL == nil || *config.URL == "" {
		return fmt.Errorf("no url found or url is empty\n")
	}
	data := []byte{}
	if d, ok := pokemonCache.Get(*config.URL); ok {
		fmt.Println("cache hit")
		data = d
	} else {
		var err error
		data, err = apiCall(config, false)
		if err != nil {
			return err
		}
		pokemonCache.Add(*config.URL, data)
	}

	var err error
	data, err = apiCall(config, true)
	if err != nil {
		return err
	}

	var locations dm.Location
	if err := json.Unmarshal(data, &locations); err != nil {
		return fmt.Errorf("error unmarshalling response body: %v\n", err)
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandExplore(config *dm.Config, pokemonCache *cache.Cache, location string) error {
	URL := "https://pokeapi.co/api/v2/location-area/" + location + "/"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v\n", err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v\n", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v\n", err)
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("error calling api: %v\n", res.Status)
	}

	var pokemons dm.LocationExplorer
	if err := json.Unmarshal(data, &pokemons); err != nil {
		return fmt.Errorf("error unmarshalling response body: %v\n", err)
	}

	if len(pokemons.PokemonEncounters) == 0 {
		fmt.Println("No Pokemons found")
	} else {
		fmt.Println("Found Pokemons:")
		for _, pokemon := range pokemons.PokemonEncounters {
			fmt.Printf("- %s\n", pokemon.Pokemon.Name)
		}
	}
	return nil
}

func CommandCatch(config *dm.Config, pokemonCache *cache.Cache, pokemon string) error {
	// Always remember to seed your random number generator
	rand.Seed(time.Now().UnixNano())

	// Get random number between 0 and 100
	randomNum := rand.Intn(200)
	URL := "https://pokeapi.co/api/v2/pokemon/" + pokemon + "/"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v\n", err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v\n", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v\n", err)
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("error calling api: %v\n", res.Status)
	}

	var pokemonData dm.Pokemon
	if err := json.Unmarshal(data, &pokemonData); err != nil {
		return fmt.Errorf("error unmarshalling response body: %v\n", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	var threshold float64
	if pokemonData.BaseExperience > 80 {
		threshold = float64(pokemonData.BaseExperience) * 0.5 // harder to catch
	} else {
		threshold = float64(pokemonData.BaseExperience) * 0.3 // easier to catch
	}

	if randomNum > int(threshold) {
		fmt.Printf("%s was caught!\n", pokemon)
		podex.AddPokemon(pokemonData)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}
	return nil
}

func CommandInspect(config *dm.Config, pokemonCache *cache.Cache, pokemon string) error {
	if _, exists := podex.Pokemon[pokemon]; exists {
		fmt.Printf("Name: %s\n", podex.Pokemon[pokemon].Name)
		fmt.Printf("Height: %d\n", podex.Pokemon[pokemon].Height)
		fmt.Printf("Weight: %d\n", podex.Pokemon[pokemon].Weight)
		fmt.Println("Stats:")
		for _, stat := range podex.Pokemon[pokemon].Stats {
			fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, type_ := range podex.Pokemon[pokemon].Types {
			fmt.Printf("- %s\n", type_.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func CommandPokedex(config *dm.Config, pokemonCache *cache.Cache, pokemon string) error {
	fmt.Println("Your pokedex:")
	for _, pokemon := range podex.Pokemon {
		fmt.Printf("- %s\n", pokemon.Name)
	}

	return nil
}
