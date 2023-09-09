package db

import "strings"

// MakeWordFreq takes a text and returns a map of words to their frequencies
func MakeWordFreq(text string) map[string]int {
	wordFreq := make(map[string]int)
	words := SplitWords(text)
	for _, word := range words {
		wordFreq[word]++
	}
	return wordFreq
}

// Split a text into words
// The delimiters used to split the text into words are defined in the function
func SplitWords(text string) (words []string) {
	// Define the delimiters used to split the text into words
	delimiters := []string{" ", "\t", "\n", ",", ".", "!", "?", ";", ":"}

	// Use the strings.FieldsFunc function to split the text into words
	fieldsFunc := func(r rune) bool {
		for _, delimiter := range delimiters {
			if string(r) == delimiter {
				return true
			}
		}
		return false
	}
	words = strings.FieldsFunc(text, fieldsFunc)

	return words
}
