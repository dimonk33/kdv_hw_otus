package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func getKeysArr(arr map[string]int) []string {
	keys := make([]string, 0, len(arr))
	for key := range arr {
		keys = append(keys, key)
	}
	return keys
}

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

	keysArr := getKeysArr(wordMap)

	sort.Slice(keysArr, func(i, j int) bool {
		if wordMap[keysArr[i]] == wordMap[keysArr[j]] {
			return keysArr[i] < keysArr[j]
		}
		return wordMap[keysArr[i]] > wordMap[keysArr[j]]
	})

	sortWordCount := len(wordMap)
	if sortWordCount > 10 {
		sortWordCount = 10
	}

	return keysArr[:sortWordCount]
}
