package main

import (
	"bufio"
	"fmt"
	"github.com/ricardosilva86/pokedexcli/cmd"
	"github.com/ricardosilva86/pokedexcli/datamodels"
	internal "github.com/ricardosilva86/pokedexcli/internal/config"
	cache "github.com/ricardosilva86/pokedexcli/internal/pokecache"
	"os"
	"strings"
	"time"
)

var (
	config = &datamodels.Config{
		URL:         &internal.BaseURL,
		PreviousURL: nil,
	}
	pokemonCache = cache.NewCache(5 * time.Minute)
)

func main() {

	go pokemonCache.ReapLoop(5 * time.Minute)

	startRepl()
}

func startRepl() {
	var input string
	var arg1 string
	for {
		fmt.Fprint(os.Stdout, "pokedex > ")
		inputReader := bufio.NewReader(os.Stdin)
		input, _ = inputReader.ReadString('\n')

		input = strings.TrimSuffix(input, "\n")
		c := strings.Split(input, " ")
		c1 := strings.ToLower(c[0])
		if len(c) >= 2 {
			arg1 = strings.ToLower(c[1])
		}

		if len(c1) == 0 {
			continue
		}
		if command, ok := cmd.GetCommands()[c1]; ok {
			err := command.Command(config, pokemonCache, arg1)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Command not found")
			continue
		}
	}
}
