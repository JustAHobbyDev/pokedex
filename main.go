package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
)

type cliCommand struct {
    name string
    description string
    callback func() error
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

func init() {
    commandExit := func() error {
        fmt.Println("Closing the Pokedex... Goodbye!")
        os.Exit(0)
        return nil
    }

    commandHelp := func() error {
        fmt.Println("Welcome to the Pokedex!")
        fmt.Println("Usage:")
        fmt.Println()

        for _, command := range commands {
            fmt.Printf("%s: %s\n", command.name, command.description)
        }

        return nil
    }

    commands = map[string]cliCommand{
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

        if err := cmd.callback(); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    }
}
