package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GroqClient struct {
	APIKey string
	Model  string
	client *http.Client
}

func NewGroqClient(apiKey, model string) *GroqClient {
	if model == "" {
		model = "meta-llama/llama-4-scout-17b-16e-instruct" // Default to Llama 3 70B
	}
	return &GroqClient{
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Message represents a chat message for Groq API
type Message struct {
	Role    string `json:"role"`    // "system", "user", or "assistant"
	Content string `json:"content"` // The actual text
}

// ChatRequest is what we send to Groq
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse is what Groq sends back
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// AskQuestion sends question + relevant chunks to Groq, returns answer
// GenerateResponse sends question with pre-built context to Groq
// This is used when context is already formatted (from ChromaDB or FTS5)
func (c *GroqClient) GenerateResponse(question, context string) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("no API key configured")
	}

	// Construct the prompt with pre-built context
	prompt := fmt.Sprintf(`You are a helpful assistant answering questions based on the provided company documents.

Documents:
%s

Question: %s

Instructions:
- Answer based ONLY on the documents provided above
- If the answer isn't in the documents, say "I don't have enough information to answer that"
- Cite which source number you used (e.g., "According to Source 1...")
- Be concise but complete`, context, question)

	// Prepare request
	reqBody := ChatRequest{
		Model: c.Model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for API errors
	if result.Error != nil {
		return "", fmt.Errorf("Groq API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return result.Choices[0].Message.Content, nil
}

// ChunkResult holds chunk data for the AI context
type ChunkResult struct {
	Content    string
	PageNumber int
	DocumentID uint
}
