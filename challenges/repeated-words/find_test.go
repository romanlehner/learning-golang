package main

import "testing"

func TestFindRepeatedWords(t *testing.T) {

	table := []struct {
		text          string
		repeatedWords string
	}{
		{text: "hello hello where in the world is this world going", repeatedWords: "hello,world"},
		{text: "hello hello where in the world is this world going hello", repeatedWords: "hello,world"},
		{text: "hello world", repeatedWords: ""},
	}

	for _, tt := range table {
		got := findRepeatedWords(tt.text)

		if got != tt.repeatedWords {
			t.Errorf("want [%v], but got [%v]", tt.repeatedWords, got)
		}
	}
}
