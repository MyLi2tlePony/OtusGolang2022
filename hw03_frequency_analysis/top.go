package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regexpForWords = regexp.MustCompile(`(\p{L}+-\p{L}+)|(\p{L}+)`)

func Top10(text string) []string {
	if text == "" {
		return make([]string, 0)
	}

	text = strings.ToLower(text)
	words := regexpForWords.FindAllString(text, -1)

	wordsFrequencies := make(map[string]int)
	for _, word := range words {
		wordsFrequencies[word]++
	}

	uniqueWords := make([]string, 0)
	for word := range wordsFrequencies {
		uniqueWords = append(uniqueWords, word)
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		if wordsFrequencies[uniqueWords[i]] == wordsFrequencies[uniqueWords[j]] {
			return uniqueWords[i] < uniqueWords[j]
		}

		return wordsFrequencies[uniqueWords[i]] > wordsFrequencies[uniqueWords[j]]
	})

	if top10 := 10; len(uniqueWords) > top10 {
		return uniqueWords[:top10]
	}

	return uniqueWords
}
