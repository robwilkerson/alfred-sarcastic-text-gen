package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

const LIMIT = 5

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

// Sarcastify takes a string and randomly converts the case of each letter to
// produce a string of mixed case letters to represent sarcasm.
//
// Example:
//
//	input:  "  Hello, World!  "
//	output: "hElLO, WOrLd!" (or some other random capitalization)
func Sarcastify(str string) string {
	normalized := normalizeString(str)

	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed for randomness.
	sarcasm := make([]rune, len(normalized))               // Use runes to handle Unicode correctly.

	for i, r := range normalized {
		if unicode.IsLetter(r) {
			if rng.Float64() < 0.5 { // 50% chance of capitalization
				sarcasm[i] = unicode.ToUpper(r)
			} else {
				sarcasm[i] = r
			}
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
		os.Exit(1) // Exit with an error code.
	}

	// Get the input string from the command line argument.
	input := os.Args[1]

	var results = []Item{}

	// Call the function to normalize and randomly capitalize the string.
	for i := 0; i < LIMIT; i++ {
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
