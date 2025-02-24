package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"strings"

	"github.com/fatih/color"
)

//go:embed words.json
var wordsJSON []byte // Embedded JSON file

type Translation struct {
	Word         string   `json:"word"`
	Translations []string `json:"translations"`
	Types        []Type   `json:"types"`
}

type Type struct {
	Noun *NounType `json:"noun,omitempty"`
	Verb *VerbType `json:"verb,omitempty"`
}

type NounType struct {
	Translation  string            `json:"translation"`
	Translations []string          `json:"translations"`
	Synonyms     map[string]string `json:"synonyms"`
	Antonyms     map[string]string `json:"antonyms"`
	Examples     map[string]string `json:"examples"`
	Frequency    int               `json:"frequency"`
}

type VerbType struct {
	Translation  string            `json:"translation"`
	Translations []string          `json:"translations"`
	Irregular    bool              `json:"irregular"`
	Synonyms     map[string]string `json:"synonyms"`
	Antonyms     map[string]string `json:"antonyms"`
	Examples     map[string]string `json:"examples"`
	Frequency    int               `json:"frequency"`
}

func main() {
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage: wordly [options]")
		fmt.Println("Options:")
		fmt.Println("  --help   Show this help message")
		return
	}

	var translations []Translation
	err := json.Unmarshal(wordsJSON, &translations)
	if err != nil {
		fmt.Println("Error parsing embedded JSON:", err)
		return
	}

	// rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(translations), func(i, j int) {
		translations[i], translations[j] = translations[j], translations[i]
	})

	for _, translation := range translations {
		fmt.Printf(" %s -> ", translation.Word)
		var userAnswer string
		fmt.Scanln(&userAnswer)

		userAnswer = strings.ToLower(strings.TrimSpace(userAnswer))
		if contains(translation.Translations, userAnswer) {
			color.Green(" Correct!\n\n")
		} else {
			color.Blue("\n Possible Translations:\n")
			showRandomTranslations(translation.Translations)

			if translation.Types != nil {
				displayWordDetails(translation.Types)
			}

			color.Blue(strings.Repeat("-", 75))
		}
	}
}

func contains(slice []string, word string) bool {
	for _, item := range slice {
		if item == word {
			return true
		}
	}
	return false
}

func showRandomTranslations(translations []string) {
	for i := 0; i < 5 && i < len(translations); i++ {
		color.Blue("  - %s", translations[rand.Intn(len(translations))])
	}
	fmt.Println()
}

func displayWordDetails(types []Type) {
	for _, t := range types {
		if t.Noun != nil {
			color.White("\n--- [ Noun ] ---")
			displayDetails(t.Noun.Synonyms, t.Noun.Antonyms, t.Noun.Examples)
		}
		if t.Verb != nil {
			color.White("\n--- [ Verb ] ---")
			displayDetails(t.Verb.Synonyms, t.Verb.Antonyms, t.Verb.Examples)
		}
	}
}

func displayDetails(synonyms, antonyms, examples map[string]string) {
	if len(synonyms) > 0 {
		color.Yellow("\n Synonyms:")
		for key, value := range synonyms {
			color.Yellow("   %s -> %s", key, value)
		}
	}
	if len(antonyms) > 0 {
		color.Red("\n Antonyms:")
		for key, value := range antonyms {
			color.Red("   %s -> %s", key, value)
		}
	}
	if len(examples) > 0 {
		color.Green("\n Examples:")
		for key, value := range examples {
			color.Green("   %s -> %s", key, value)
		}
	}
	fmt.Println()
}
