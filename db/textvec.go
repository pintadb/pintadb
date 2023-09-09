package db

import (
	"strings"
)

type TextVec struct {
	RawText  string
	WordFreq map[string]int
	Vector   []float64
}

// CalculateVector takes a text and calculates its vector representation
// The text is split into words using the SplitWords functionthen each the
// word frequencies are then used to calculate the vector representation of the text
func (t *TextVec) CalculateVector(dimension uint64) {
	totalWords := 0
	for _, v := range t.WordFreq {
		totalWords += v
	}

	for k, v := range t.WordFreq {
		weight := float64(v) / float64(totalWords)
		vectorValue := Text2Vec(k, dimension)
		for i, vv := range vectorValue {
			t.Vector[i] += weight * vv
		}
	}
}

// Text2Vec converts a text to a vector of floats
// A vector is initialized with zeros and then updated
// with word counts that are hashed to a random index in the vector
// After that the vector is returned
// The dimension should be larger than the number of unique words in the text
func Text2Vec(text string, dimension uint64) []float64 {
	// Convert the input text to lowercase and split it into words
	words := strings.Split(strings.ToLower(text), " ")

	// Initialize a map to hold the word counts
	wordCounts := make(map[string]int)

	// Count the number of occurrences of each word
	for _, word := range words {
		wordCounts[word]++
	}

	// Initialize a vector of zeros with the specified dimension
	vector := make([]float64, dimension)

	// Update the vector with the word counts
	for word, count := range wordCounts {
		// Hash the word to an integer between 0 and dimension-1
		hash := uint32(FowlerNollVo32(word)) % uint32(dimension)

		// Use the hash as the index into the vector
		vector[hash] += float64(count)
	}

	return vector
}
