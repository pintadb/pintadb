package db

import (
	"fmt"
	"strings"
)

type TextVec struct {
	RawText  string
	WordFreq map[string]int
	Vector   []float64
}

func (t *TextVec) PrintRawText() {
	fmt.Println(t.RawText)
}

// CalculateVector takes a text and calculates its vector representation
// The text is split into words using the SplitWords function
// The word frequencies are calculated using the map[string]int type
// The word frequencies are then used to calculate the vector representation of the text
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

// Convert a text to a vector of floats
// The vector is initialized with zeros and then updated with the word counts
// The word counts are hashed to a random index in the vector
// The vector is returned
// The vector is used to represent the text in a machine learning model
// The dimension is a hyperparameter that can be tuned
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
		hash := uint32(fowlerNollVo32(word)) % uint32(dimension)

		// Use the hash as the index into the vector
		vector[hash] += float64(count)
	}

	return vector
}
