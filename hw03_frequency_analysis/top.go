package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCount struct {
	word  string
	count int
}

func Top10(str string) []string {
	wordsCount := getWordsCount(getWordsCountMap(str))

	sort.Slice(wordsCount, func(i, j int) bool {
		if wordsCount[i].count == wordsCount[j].count {
			return wordsCount[i].word < wordsCount[j].word
		}

		return wordsCount[i].count > wordsCount[j].count
	})

	return getTopWords(wordsCount)
}

func getTopWords(wordsCount []WordCount) []string {
	maxLength := 10
	wordsCountLength := len(wordsCount)
	var topWords []string

	if wordsCountLength < maxLength {
		maxLength = wordsCountLength
	}

	for i := 0; i < maxLength; i++ {
		topWords = append(topWords, wordsCount[i].word)
	}

	return topWords
}

func getWordsCount(wordsCountMap map[string]int) []WordCount {
	wordsCount := make([]WordCount, 0)

	for word, count := range wordsCountMap {
		wordsCount = append(wordsCount, WordCount{word, count})
	}

	return wordsCount
}

func getWordsCountMap(str string) map[string]int {
	wordsCountMap := make(map[string]int)

	words := strings.Fields(str)

	for _, word := range words {
		wordsCountMap[word]++
	}

	return wordsCountMap
}
