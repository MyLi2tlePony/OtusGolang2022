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

	sortWords(&words)

	wordsFrequency := make(map[string]int)

	for _, word := range words {
		wordsFrequency[word]++
	}

	top10 := 10
	topFrequencies := findTopFrequencies(wordsFrequency, top10)
	wordsWithTopFrequency := findWordsWithTopFrequency(wordsFrequency, topFrequencies)

	return selectTopWordsWithTopFrequency(wordsWithTopFrequency, top10)
}

func sortWords(unsortedWords *[]string) {
	sort.Slice(*unsortedWords, func(i, j int) bool {
		return (*unsortedWords)[i] < (*unsortedWords)[j]
	})
}

func findTopFrequencies(wordsFrequency map[string]int, topFrequenciesNumber int) []int {
	topFrequencies := make([]int, topFrequenciesNumber)

	for _, frequency := range wordsFrequency {
		if frequency > topFrequencies[0] {
			topFrequencies[0] = frequency

			sort.Slice(topFrequencies, func(i, j int) bool {
				return topFrequencies[i] < topFrequencies[j]
			})
		}
	}

	return topFrequencies
}

func findWordsWithTopFrequency(wordsFrequency map[string]int, topFrequencies []int) [][]string {
	sort.Slice(topFrequencies, func(i, j int) bool {
		return topFrequencies[i] > topFrequencies[j]
	})

	wordsWithTopFrequency := make([][]string, len(topFrequencies))

	for word, frequency := range wordsFrequency {
		for index, value := range topFrequencies {
			if frequency == value {
				wordsWithTopFrequency[index] = append(wordsWithTopFrequency[index], word)
				break
			}
		}
	}

	return wordsWithTopFrequency
}

func selectTopWordsWithTopFrequency(wordsWithTopFrequency [][]string, topWordsNumber int) []string {
	topWordsWithTopFrequency := make([]string, 0)

	for wordIndex := 0; wordIndex < len(wordsWithTopFrequency); wordIndex++ {
		sortWords(&wordsWithTopFrequency[wordIndex])

		for _, word := range wordsWithTopFrequency[wordIndex] {
			topWordsWithTopFrequency = append(topWordsWithTopFrequency, word)

			if len(topWordsWithTopFrequency) >= topWordsNumber {
				return topWordsWithTopFrequency
			}
		}
	}

	return topWordsWithTopFrequency
}
