package db

import "math"

// CosineDistance calculates the cosine distance between two vectors
func CosineDistance(a, b []float64) float64 {
	// Calculate the dot product of the two vectors
	dotProduct := 0.0
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
	}

	// Calculate the magnitudes of the two vectors
	magnitudeA := 0.0
	for _, val := range a {
		magnitudeA += val * val
	}
	magnitudeA = math.Sqrt(magnitudeA)

	magnitudeB := 0.0
	for _, val := range b {
		magnitudeB += val * val
	}
	magnitudeB = math.Sqrt(magnitudeB)

	// Calculate the cosine similarity
	cosineSimilarity := dotProduct / (magnitudeA * magnitudeB)

	// Calculate the cosine distance
	cosineDistance := 1 - cosineSimilarity

	return cosineDistance
}
