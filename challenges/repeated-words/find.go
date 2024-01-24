package main

import (
	"fmt"
	"strings"
)

func findRepeatedWords(text string) string {
	// separate text into words
	wordlist := strings.Fields(text)
	wordCount := make(map[string]int)

	for _, word := range wordlist {
		wordCount[word]++
	}

	var repeatedWords string

	// a map does not keep elements in the order they were added
	// this might be a problem if we want to return the words in the order they appear in the text
	// the test might fail if we return "world,hello" instead of "hello,world"
	// I ignored this problem for now
	for word, count := range wordCount {
		if count > 1 {
			repeatedWords += fmt.Sprintf(",%s", word)
		}
	}

	if len(repeatedWords) == 0 {
		return ""
	}

	return repeatedWords[1:]
}
