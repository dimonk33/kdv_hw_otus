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
		_, ok := wordMap[word]
		if !ok {
			wordMap[word] = 0
		}
		wordMap[word]++
	}

	type keyValue struct {
		key   string
		count int
	}
	sortedMap := make([]keyValue, len(wordMap))
	i := 0
	for k, v := range wordMap {
		sortedMap[i] = keyValue{key: k, count: v}
		i++
	}

	sort.Slice(sortedMap, func(i, j int) bool {
		if sortedMap[i].count == sortedMap[j].count {
			return sortedMap[i].key < sortedMap[j].key
		}
		return sortedMap[i].count > sortedMap[j].count
	})

	sortWords := make([]string, 0)
	for _, v := range sortedMap {
		sortWords = append(sortWords, v.key)
		if len(sortWords) >= 10 {
			break
		}
	}

	return sortWords
}
