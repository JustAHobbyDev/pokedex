package main

import (
    "fmt"
    "strings"
)

func cleanInput(text string) []string {
    splits := strings.Split(text, " ")
    words := []string{}

    for _, split := range splits {
        if len(split) > 0 {
            words = append(words, split)
        }
    }

    return words
}

func main() {
    fmt.Println("Hello, World!")
}
