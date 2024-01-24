package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("Add a word and definition", func(t *testing.T) {
		dictionary := Dictionary{}
		word := "test"
		definition := "this is just a test"

		err := dictionary.Add(word, definition)

		if err != nil {
			t.Errorf("got error: %s, but want error: %v", err, nil)
		}
	})

	t.Run("Add existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}

		err := dictionary.Add(word, definition)

		if err != ErrorWordExists {
			t.Errorf("got error: %s, but want error: %s", err, ErrorWordExists)
		}
	})
}

func TestSearch(t *testing.T) {
	t.Run("Search for a word that doesn't exist", func(t *testing.T) {
		dictionary := Dictionary{}
		got, err := dictionary.Search("test")
		want := ""

		if got != want {
			t.Errorf("got: %s, but want an empty string >%s<", got, want)
		}

		if err != ErrorNotFound {
			t.Errorf("got error: %s, but want error: %s", err, ErrorNotFound)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update definition", func(t *testing.T) {
		dictionary := Dictionary{"test": "this is just a test"}
		word := "test"
		want := "this is serious now"

		err := dictionary.Update(word, want)

		if err != nil {
			t.Errorf("got error: %s, but want error: %v", err, nil)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete a word from the dictionary", func(t *testing.T) {
		dictionary := Dictionary{"test": "this is just a test"}
		word := "test"

		dictionary.Delete(word)

		_, err := dictionary.Search(word)

		if err != ErrorNotFound {
			t.Errorf("got error: %s, but want error: %s", err, ErrorNotFound)
		}
	})
}