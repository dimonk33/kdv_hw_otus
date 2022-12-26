package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	wordMap := map[string]int{}
	wordArr := strings.FieldsFunc(text, func(r rune) bool {
		return strings.ContainsRune(" ,.:!?\n\t\"'`", r)
	})
	for _, word := range wordArr {
		word = strings.ToLower(strings.TrimSpace(word))
		if word == "" || word == "-" {
			continue
		}
		wordMap[word]++
	}

	type keyValue struct {
		key   string
		count int
	}
	sortedMap := make([]keyValue, 0, len(wordMap))
	for k, v := range wordMap {
		sortedMap = append(sortedMap, keyValue{key: k, count: v})
	}

	sort.Slice(sortedMap, func(i, j int) bool {
		if sortedMap[i].count == sortedMap[j].count {
			return sortedMap[i].key < sortedMap[j].key
		}
		return sortedMap[i].count > sortedMap[j].count
	})

	var sortWordCount int
	if len(sortedMap) > 10 {
		sortWordCount = 10
	} else {
		sortWordCount = len(sortedMap)
	}

	sortWords := make([]string, sortWordCount)
	for i := 0; i < sortWordCount; i++ {
		sortWords[i] = sortedMap[i].key
	}

	return sortWords
}
