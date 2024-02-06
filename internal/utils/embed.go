package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pgvector/pgvector-go"
)

// Embed a string
func Embed(input string) (pgvector.Vector, error) {
	// Model API url
	apiUrl := "http://llm:8080/embed"

	// Construct payload
	payload := []string{input}

	// Serialize the payload to JSON
	payloadBytes, err := json.Marshal(payload)

	// Error check
	if err != nil {
		return pgvector.Vector{}, err
	}

	// Make a POST request to the FastAPI /embed route
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(payloadBytes))

	// Make sure to close body
	defer resp.Body.Close()

	// Error check
	if err != nil {
		return pgvector.Vector{}, err
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return pgvector.Vector{}, err
	}

	// Create decoder
	decoder := json.NewDecoder(resp.Body)

	// Create output embedding
	var embeddings [][]float32

	// Decode embedding
	err = decoder.Decode(&embeddings)

	// Error check
	if err != nil {
		return pgvector.Vector{}, err
	}

	// Convert the zero array to a Vector
	return pgvector.NewVector(embeddings[0]), nil
}
