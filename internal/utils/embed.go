package utils

import (
	"github.com/pgvector/pgvector-go"
)

// Embed a string
func Embed(input string) pgvector.Vector {
	// Create a 1024-zero array
	zeroVector := make([]float32, 1024)

	// Convert the zero array to a Vector
	return pgvector.NewVector(zeroVector)
}
