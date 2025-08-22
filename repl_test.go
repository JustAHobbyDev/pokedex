package main

import (
    "testing"
)

func TestCleanInput(t *testing.T) {
    cases := []struct {
        input string
        expected []string
    }{
        {
            input: " hello world ",
            expected: []string{ "hello", "world"},
        },
        {
            input: "hello wide world ",
            expected: []string{ "hello", "wide", "world"},
        },
        {
            input: " ",
            expected: []string{},
        },
    }

    for _, c := range cases {
        actual := cleanInput(c.input)
        if len(actual) != len(c.expected) {
            t.Errorf("Lengths do not match\nactual: %d\nexpected: %d",
                len(actual), len(c.expected))
        }

        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]

            if word != expectedWord {
                t.Errorf("FAIL:\n\tactual: %s\n\texpected: %s", word, expectedWord)
            }
        }
    }
}
