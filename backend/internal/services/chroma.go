package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ChromaClient is a simple HTTP client for ChromaDB
type ChromaClient struct {
	baseURL      string
	httpClient   *http.Client
	collection   string
	enabled      bool
	collectionID string // Required UUID for v2 data operations
}

// NewChromaClient creates a new ChromaDB HTTP client
func NewChromaClient(baseURL string) *ChromaClient {
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	fmt.Printf("Initializing ChromaDB client with URL: %s\n", baseURL)

	client := &ChromaClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		collection: "enterprise_brain",
		enabled:    false,
	}

	// Test connection and create/find collection
	if err := client.initialize(); err != nil {
		fmt.Printf("ChromaDB initialization failed: %v\n", err)
		return nil
	}

	client.enabled = true
	fmt.Printf("ChromaDB client initialized successfully\n")
	return client
}

func (c *ChromaClient) initialize() error {
	// Step 1: Heartbeat (v2)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v2/heartbeat", c.baseURL), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("heartbeat failed: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("heartbeat status: %d", resp.StatusCode)
	}
	fmt.Printf("ChromaDB heartbeat successful\n")

	// Step 2: List collections to find our collection
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections", c.baseURL)
	req, err = http.NewRequestWithContext(ctx2, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("list collections failed status: %d", resp.StatusCode)
	}

	// Parse the response - v2 returns a list of collections
	var collections []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&collections); err != nil {
		return fmt.Errorf("failed to decode collections: %w", err)
	}

	// Find our collection by name
	for _, col := range collections {
		if col.Name == c.collection {
			c.collectionID = col.ID
			fmt.Printf("Collection '%s' exists with ID: %s\n", c.collection, c.collectionID)
			return nil
		}
	}

	// Collection not found, create it
	fmt.Printf("Collection '%s' not found, creating...\n", c.collection)
	return c.createCollection()
}

func (c *ChromaClient) createCollection() error {
	body := map[string]interface{}{"name": c.collection}
	jsonBody, _ := json.Marshal(body)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		resBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create failed: status=%d, body=%s", resp.StatusCode, string(resBody))
	}

	var colObj struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&colObj); err != nil {
		return fmt.Errorf("failed to decode created collection: %w", err)
	}
	c.collectionID = colObj.ID
	fmt.Printf("Collection '%s' created successfully with ID: %s\n", c.collection, c.collectionID)
	return nil
}

func (c *ChromaClient) IsHealthy() bool {
	if c == nil || !c.enabled {
		return false
	}
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v2/heartbeat", c.baseURL))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (c *ChromaClient) StoreChunks(chunks []ChunkWithEmbedding) error {
	if !c.IsHealthy() {
		return fmt.Errorf("chroma client not healthy")
	}

	var ids []string
	var documents []string
	var embeddings [][]float32
	var metadatas []map[string]interface{}

	for _, chunk := range chunks {
		ids = append(ids, fmt.Sprintf("chunk_%d", chunk.ID))
		documents = append(documents, chunk.Content)
		embeddings = append(embeddings, chunk.Embedding)
		metadatas = append(metadatas, map[string]interface{}{
			"document_id": chunk.DocumentID,
			"page_number": chunk.PageNumber,
			"chunk_index": chunk.ChunkIndex,
		})
	}

	body := map[string]interface{}{
		"ids": ids, "documents": documents, "embeddings": embeddings, "metadatas": metadatas,
	}
	jsonBody, _ := json.Marshal(body)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use v2 path with UUID
	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s/add", c.baseURL, c.collectionID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		resBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add chunks: %s", string(resBody))
	}
	return nil
}

func (c *ChromaClient) Search(queryEmbedding []float32, topK int) ([]SearchResult, error) {
	if !c.IsHealthy() {
		return nil, fmt.Errorf("chroma client not healthy")
	}

	body := map[string]interface{}{
		"query_embeddings": [][]float32{queryEmbedding},
		"n_results":        topK,
	}
	jsonBody, _ := json.Marshal(body)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use v2 path with UUID
	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s/query", c.baseURL, c.collectionID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		IDs       [][]string                 `json:"ids"`
		Documents [][]string                 `json:"documents"`
		Metadatas [][]map[string]interface{} `json:"metadatas"`
		Distances [][]float64                `json:"distances"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var searchResults []SearchResult
	if len(result.Documents) > 0 && len(result.Documents[0]) > 0 {
		for i, doc := range result.Documents[0] {
			meta := result.Metadatas[0][i]
			searchResults = append(searchResults, SearchResult{
				Content:    doc,
				DocumentID: uint(meta["document_id"].(float64)),
				PageNumber: int(meta["page_number"].(float64)),
				Score:      result.Distances[0][i],
			})
		}
	}
	return searchResults, nil
}

func (c *ChromaClient) DeleteChunksByDocumentID(chunkIDs []uint) error {
	if !c.IsHealthy() || len(chunkIDs) == 0 {
		return nil
	}

	var ids []string
	for _, id := range chunkIDs {
		ids = append(ids, fmt.Sprintf("chunk_%d", id))
	}

	body := map[string]interface{}{"ids": ids}
	jsonBody, _ := json.Marshal(body)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use v2 path with UUID and /delete endpoint
	url := fmt.Sprintf("%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s/delete", c.baseURL, c.collectionID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		resBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(resBody))
	}
	return nil
}

// ChunkWithEmbedding represents a chunk ready for vector storage
type ChunkWithEmbedding struct {
	ID         uint
	DocumentID uint
	Content    string
	PageNumber int
	ChunkIndex int
	Embedding  []float32
}

// SearchResult represents a found chunk from vector search
type SearchResult struct {
	Content    string
	DocumentID uint
	PageNumber int
	Score      float64
}
