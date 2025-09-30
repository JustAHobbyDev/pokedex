package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "io"
    "net/http"
    "encoding/json"
    "github.com/JustAHobbyDev/pokedex/internal/pokecache"
)

type Config struct {
    Next     string
    Previous string
}

type cliCommand struct {
    name        string
    description string
    callback    func(c *Config) error
}

var commands map[string]cliCommand

func cleanInput(text string) []string {
    splits := strings.Split(text, " ")
    words := []string{}

    for _, split := range splits {
        if len(split) > 0 {
            words = append(words, strings.ToLower(split))
        }
    }

    return words
}

type LocationArea struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}

type LocationAreas struct {
    Count    int            `json:"count"`
    Next     *string        `json:"next"`
    Previous *string        `json:"previous"`
    Results  []LocationArea `json:"results"`
}

var cache pokecache.Cache;

func fetchOrGet(url string) ([]byte, error) {
    cachedRes, ok := cache.Get(url)
    if ok {
        return cachedRes, nil
    }

    res, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("Failed to GET %s [Error: %v]\n", url, err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, fmt.Errorf("Error reading response body: %v\n", err)
    }

    cache.Add(url, body)

    return body, nil
}

func init() {
    cache = pokecache.NewCache(5)

    commandMap := func(c *Config) error {
        if c.Next == "" {
            fmt.Println("you're on the first page")
        }

        data, err := fetchOrGet(c.Next)

        if err != nil {
            fmt.Println(err)
            return nil
        }

        var locationAreas LocationAreas
        if err := json.Unmarshal(data, &locationAreas); err != nil {
            fmt.Println("Error parsing JSON: ", err)
            return nil
        }

        results := locationAreas.Results
        for _, result := range results {
            fmt.Println(result.Name)
        }

        c.Previous = c.Next
        c.Next = *locationAreas.Next

        return nil
    }

    commandMapb := func(c *Config) error {
        if c.Previous == "" {
            fmt.Println("you're on the first page")
            c.Previous = c.Next
        }

        data, err := fetchOrGet(c.Previous)
        if err != nil {
            fmt.Println(err)
            return nil
        }

        var locationAreas LocationAreas
        if err := json.Unmarshal(data, &locationAreas); err != nil {
            fmt.Println("Error parsing JSON: ", err)
            return nil
        }

        results := locationAreas.Results
        for _, result := range results {
            fmt.Println(result.Name)
        }

        if locationAreas.Previous != nil {
            c.Previous = *locationAreas.Previous
        } else {
            c.Previous = ""
        }
        c.Next = *locationAreas.Next

        return nil
    }

    commandExit := func(c *Config) error {
        fmt.Println("Closing the Pokedex... Goodbye!")
        os.Exit(0)
        return nil
    }

    commandHelp := func(c *Config) error {
        fmt.Println("Welcome to the Pokedex!")
        fmt.Println("Usage:")
        fmt.Println()

        for _, command := range commands {
            fmt.Printf("%s: %s\n", command.name, command.description)
        }

        return nil
    }

    commands = map[string]cliCommand{
        "map": {
            name: "map",
            description: "Page forward through locations",
            callback: commandMap,
        },
        "mapb": {
            name: "mapb",
            description: "Page backward through locations",
            callback: commandMapb,
        },
        "exit": {
            name: "exit",
            description: "Exit the Pokedex",
            callback: commandExit,
        },
        "help": {
            name: "help",
            description: "Print Pokedex help",
            callback: commandHelp,
        },
    }
}

func main() {
    cache.ReapLoop()

    // https://pokeapi.co/api/v2/location/{id or name}/
    config := &Config{
        Next: "https://pokeapi.co/api/v2/location-area/",
    }

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex > ")
        if !scanner.Scan() {
            if err := scanner.Err(); err != nil {
                fmt.Fprintln(os.Stderr, "reading standard input:", err)
            }
            break
        }

        inputWords := cleanInput(scanner.Text())
        if len(inputWords) == 0 {
            fmt.Println("Please enter command...")
            continue
        }

        input := inputWords[0]
        cmd, exists := commands[input]
        if !exists {
            fmt.Printf("Unknown command: %s\n", input)
            continue
        }

        if err := cmd.callback(config); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    }
}
