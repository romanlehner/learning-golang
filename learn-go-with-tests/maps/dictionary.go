package main

const (
	ErrorNotFound   = DictionaryErr("could not find the word you were looking for")
	ErrorWordExists = DictionaryErr("word already exists in the dictionary")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

type Dictionary map[string]string

func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrorNotFound:
		d[word] = definition
	case nil:
		return ErrorWordExists
	default:
		return err
	}

	return nil
}

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]

	if !ok {
		return "", ErrorNotFound
	}

	return definition, nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	if err != nil {
		return err
	}

	d[word] = definition

	return nil
}

func (d Dictionary) Delete(word string) {
	delete(d, word)
}
