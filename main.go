package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
)

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

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan() 
        input :=  scanner.Text()
        cleanInput := cleanInput(input)
        fmt.Printf("Your command was: %s\n", cleanInput[0])
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}
