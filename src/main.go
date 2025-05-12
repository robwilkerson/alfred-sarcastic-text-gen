package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Item struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
}

type AlfredOutput struct {
	Items []Item `json:"items"`
}

// normalizeString takes a string and normalizes it to...
//   - Remove leading and trailing spaces
//   - Convert to lowercase
//
// Example:
//
//	input:  " Hello, World!       "
//	output: "hello, world!"
func normalizeString(str string) string {
	if str == "" {
		return ""
	}

	normalized := strings.TrimSpace(str)
	normalized = strings.ToLower(normalized)

	return normalized
}

// disambiguateChar accepts a rune and, if the rune is one that can be ambiguous
// when displayed (e.g., lowercase "L" looks like and uppercase "I"), returns
// the disambiguated form of that rune.
//
// The function doesn't just remove a letter because it's  ambiguous, nor does
// it change a zero (0) to something else entirely go avoid the ambiguity
// because the meaning of the larger string should not be altered. It only maps
// characters to a less ambiguous character with the same value.
func disambiguateChar(r rune) rune {
	ambiguousChars := map[string]rune{
		"I": rune('i'), // Can look like lowercase "L"
		"l": rune('L'), // Can look like uppercase "I" or even like "1"
		"O": rune('o'), // Can look like zero (0)
	}

	// If the passed character is in the ambiguousChars map, replace it with its
	// unambiguous version
	if alt, ok := ambiguousChars[string(r)]; ok {
		return alt
	}

	return r
}

// Sarcastify takes a string and randomly converts the case of each letter to
// produce a string of mixed case letters to represent sarcasm.
//
// Example:
//
//	input:  "  Hello, World!  "
//	output: "hElLO, WOrLd!" (or some other random capitalization)
func Sarcastify(str string) string {
	normalized := normalizeString(str)

	rand := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed for randomness.
	sarcasm := make([]rune, len(normalized))                // Use runes to handle Unicode correctly.

	for i, r := range normalized {
		if unicode.IsLetter(r) {
			if rand.Float64() < 0.5 { // 50% chance of capitalization
				sarcasm[i] = unicode.ToUpper(r)
			} else {
				sarcasm[i] = r
			}

			// Lastly, apply any disambiguation
			sarcasm[i] = disambiguateChar(sarcasm[i])
		} else {
			sarcasm[i] = r // Keep non-letters as they are.
		}
	}

	return string(sarcasm)
}

func main() {
	// Check if an argument was provided.
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <string>")
		fmt.Println("Alfred Usage: sc <string>")
		os.Exit(1) // Exit with an error code.
	}

	// Get the input string from the command line argument.
	input := os.Args[1]

	// The options we'll return
	var results = []Item{}
	// The max number of options to return
	var limit, _ = strconv.Atoi(os.Getenv("LIMIT"))

	// Call the function to normalize and randomly capitalize the string.
	for _ = range limit {
		result := Sarcastify(input)
		results = append(results, Item{
			Title: result, // displayed
			Arg:   result, // sent to clipboard
		})
	}

	output := AlfredOutput{Items: results}
	jsonOutput, err := json.Marshal(output)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(strings.TrimSpace(string(jsonOutput)))
}
