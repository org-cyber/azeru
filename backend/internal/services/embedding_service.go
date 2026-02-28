package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
	client *api.Client
	url    string
	model  string
}

// NewOllamaClient creates a new Ollama client with valid HTTP transport
func NewOllamaClient(baseURL string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	// Parse URL first - if invalid, return nil (caller checks this)
	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("Invalid Ollama URL %s: %v\n", baseURL, err)
		return nil
	}

	// CRITICAL FIX: Always create HTTP client, never pass nil
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	// Create Ollama client with valid HTTP client
	client := api.NewClient(u, httpClient)

	return &OllamaClient{
		client: client,
		url:    baseURL,
		model:  "nomic-embed-text",
	}
}

// IsHealthy checks if Ollama is actually reachable and working
func (o *OllamaClient) IsHealthy() bool {
	if o == nil || o.client == nil {
		return false
	}

	// Try to list models - lightweight health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.client.List(ctx)
	return err == nil
}

func (o *OllamaClient) GetEmbedding(text string) ([]float32, error) {
	if o == nil || o.client == nil {
		return nil, fmt.Errorf("ollama client not initialized")
	}

	req := &api.EmbeddingRequest{
		Model:  o.model,
		Prompt: text,
	}

	resp, err := o.client.Embeddings(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// Manual conversion from []float64 to []float32
	embedding := make([]float32, len(resp.Embedding))
	for i, v := range resp.Embedding {
		embedding[i] = float32(v)
	}

	return embedding, nil
}

func (o *OllamaClient) GetEmbeddingsBatch(texts []string) ([][]float32, error) {
	if o == nil || o.client == nil {
		return nil, fmt.Errorf("ollama client not initialized")
	}

	var embeddings [][]float32
	for i, text := range texts {
		embedding, err := o.GetEmbedding(text)
		if err != nil {
			return nil, fmt.Errorf("failed to embed chunk %d: %w", i, err)
		}
		embeddings = append(embeddings, embedding)
	}
	return embeddings, nil
}
