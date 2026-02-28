package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"azeru/internal/config"
	"azeru/internal/database"
	"azeru/internal/models"
	"azeru/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	config *config.Config
	db     *gorm.DB
	ai     *services.GroqClient
	ollama *services.OllamaClient
	chroma *services.ChromaClient
	pdf    *services.PDFService
}

// NewHandler creates handler with all dependencies
func NewHandler(cfg *config.Config, db *gorm.DB, ollama *services.OllamaClient, chroma *services.ChromaClient) *Handler {
	return &Handler{
		config: cfg,
		db:     db,
		ai:     services.NewGroqClient(cfg.GroqAPIKey, cfg.GroqModel),
		ollama: ollama,
		chroma: chroma,
		pdf:    services.NewPDFService(),
	}
}

// UploadDocument handles PDF uploads
func (h *Handler) UploadDocument(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files allowed"})
		return
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	doc := models.Document{
		FileName:    header.Filename,
		FileSize:    header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Status:      "processing",
	}

	if err := h.db.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	filePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", doc.ID, doc.FileName))
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	go h.processPDF(doc.ID, filePath)

	c.JSON(http.StatusAccepted, gin.H{
		"message":  "Upload successful, processing started",
		"document": doc,
	})
}

func (h *Handler) processPDF(docID uint, filePath string) {
	var doc models.Document
	if err := h.db.First(&doc, docID).Error; err != nil {
		log.Printf("PDF %d: Document not found: %v", docID, err)
		return
	}

	// Extract text and chunk in one operation
	chunkContents := h.pdf.ChunkText(filePath, 800, 100)

	if len(chunkContents) == 0 {
		log.Printf("PDF %d: No text extracted (possibly scanned image PDF)", docID)
		doc.Status = "failed"
		h.db.Save(&doc)
		return
	}

	log.Printf("PDF %d: Extracted %d chunks", docID, len(chunkContents))

	var chunks []models.Chunk
	var chunksWithEmbeddings []services.ChunkWithEmbedding

	for i, content := range chunkContents {
		chunk := models.Chunk{
			DocumentID: docID,
			Content:    content,
			ChunkIndex: i,
			PageNum:    0,
		}
		chunks = append(chunks, chunk)
	}

	// Try to generate embeddings if Ollama is healthy
	if h.ollama != nil && h.ollama.IsHealthy() && len(chunks) > 0 {
		log.Printf("PDF %d: Generating embeddings for %d chunks...", docID, len(chunks))

		var texts []string
		for _, chunk := range chunks {
			texts = append(texts, chunk.Content)
		}

		embeddings, err := h.ollama.GetEmbeddingsBatch(texts)
		if err != nil {
			log.Printf("PDF %d: WARNING - Embedding generation failed: %v", docID, err)
			log.Printf("PDF %d: Document will be indexed without embeddings (FTS5 only)", docID)
		} else {
			log.Printf("PDF %d: Successfully generated %d embeddings", docID, len(embeddings))

			// Prepare chunks with embeddings for ChromaDB
			for i := range chunks {
				if i < len(embeddings) {
					// Store in model for SQLite
					embeddingJSON, _ := json.Marshal(embeddings[i])
					chunks[i].Embedding = embeddingJSON

					// Prepare for ChromaDB
					chunksWithEmbeddings = append(chunksWithEmbeddings, services.ChunkWithEmbedding{
						ID:         chunks[i].ID, // Will be set after DB insert
						DocumentID: docID,
						Content:    chunks[i].Content,
						PageNumber: chunks[i].PageNum,
						ChunkIndex: chunks[i].ChunkIndex,
						Embedding:  embeddings[i],
					})
				}
			}
		}
	} else {
		if h.ollama == nil {
			log.Printf("PDF %d: Ollama client not initialized - indexing without embeddings", docID)
		} else if !h.ollama.IsHealthy() {
			log.Printf("PDF %d: Ollama not reachable - indexing without embeddings", docID)
		}
	}

	// Save chunks to SQLite (get IDs assigned)
	if err := h.db.CreateInBatches(&chunks, 100).Error; err != nil {
		log.Printf("PDF %d: Failed to save chunks: %v", docID, err)
		doc.Status = "failed"
		h.db.Save(&doc)
		return
	}

	// Update IDs in chunksWithEmbeddings and save to ChromaDB
	if h.chroma != nil && h.chroma.IsHealthy() && len(chunksWithEmbeddings) > 0 {
		for i := range chunksWithEmbeddings {
			chunksWithEmbeddings[i].ID = chunks[i].ID
		}

		if err := h.chroma.StoreChunks(chunksWithEmbeddings); err != nil {
			log.Printf("PDF %d: WARNING - Failed to store in ChromaDB: %v", docID, err)
		} else {
			log.Printf("PDF %d: Stored %d chunks in ChromaDB", docID, len(chunksWithEmbeddings))
		}
	}

	doc.ChunkCount = len(chunks)
	doc.Status = "indexed"
	h.db.Save(&doc)

	log.Printf("PDF %d: Processing complete - %d chunks indexed", docID, len(chunks))
}

func (h *Handler) Chat(c *gin.Context) {
	var req struct {
		Question string `json:"question" binding:"required"`
		APIKey   string `json:"api_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question and API key required"})
		return
	}

	// Create AI client with user's key (BYOK)
	userAI := services.NewGroqClient(req.APIKey, h.config.GroqModel)

	// Search strategy: Try ChromaDB first (vector search), fallback to FTS5
	var chunks []models.Chunk
	var searchMethod string = "none"
	var searchResults []services.SearchResult
	var err error

	// Detailed logging for debugging
	log.Printf("Chat: Starting search for question: %q", req.Question)
	log.Printf("Chat: ChromaDB healthy=%v, Ollama healthy=%v", h.chroma != nil && h.chroma.IsHealthy(), h.ollama != nil && h.ollama.IsHealthy())

	if h.chroma != nil && h.chroma.IsHealthy() && h.ollama != nil && h.ollama.IsHealthy() {
		log.Printf("Chat: Attempting ChromaDB vector search...")
		// Generate embedding for query
		queryEmbedding, err := h.ollama.GetEmbedding(req.Question)
		if err != nil {
			log.Printf("Chat: ERROR generating query embedding: %v", err)
		} else {
			log.Printf("Chat: Generated query embedding, dimension=%d", len(queryEmbedding))
			// Vector search in ChromaDB
			searchResults, err = h.chroma.Search(queryEmbedding, 5)
			if err != nil {
				log.Printf("Chat: ERROR in ChromaDB search: %v", err)
			} else if len(searchResults) > 0 {
				log.Printf("Chat: SUCCESS - Found %d results via ChromaDB vector search", len(searchResults))
				searchMethod = "vector_search"
				// Convert SearchResult to models.Chunk for consistent handling
				for _, result := range searchResults {
					chunks = append(chunks, models.Chunk{
						DocumentID: result.DocumentID,
						Content:    result.Content,
						PageNum:    result.PageNumber,
					})
				}
			} else {
				log.Printf("Chat: ChromaDB returned 0 results")
			}
		}
	} else {
		log.Printf("Chat: ChromaDB/Ollama not available, skipping vector search")
	}

	// Fallback to FTS5 if ChromaDB failed or not available
	if len(chunks) == 0 {
		log.Printf("Chat: Falling back to FTS5 full-text search")
		searchMethod = "keyword_search"
		chunks, err = database.SearchChunks(h.db, req.Question, 5)
		if err != nil {
			log.Printf("Chat: ERROR in FTS5 search: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed: " + err.Error()})
			return
		}
		if len(chunks) > 0 {
			log.Printf("Chat: FTS5 search returned %d results", len(chunks))
		} else {
			log.Printf("Chat: FTS5 search returned 0 results")
		}
	}

	if len(chunks) == 0 {
		log.Printf("Chat: No relevant documents found")
		c.JSON(http.StatusOK, gin.H{
			"answer":        "I couldn't find any relevant information in the uploaded documents.",
			"sources":       []string{},
			"search_method": searchMethod,
		})
		return
	}

	// Build context from chunks
	context := h.buildContext(chunks)

	// Get answer from AI
	answer, err := userAI.GenerateResponse(req.Question, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI request failed: " + err.Error()})
		return
	}

	// Save chat messages
	h.saveChatMessage("user", req.Question)
	h.saveChatMessage("assistant", answer)

	log.Printf("Chat: Successfully returned answer using %s", searchMethod)

	c.JSON(http.StatusOK, gin.H{
		"answer":        answer,
		"sources":       chunks,
		"search_method": searchMethod,
	})
}

func (h *Handler) buildContext(chunks []models.Chunk) string {
	var context string
	for i, chunk := range chunks {
		context += fmt.Sprintf("\n[Excerpt %d - Document %d]: %s\n", i+1, chunk.DocumentID, chunk.Content)
	}
	return context
}

func (h *Handler) saveChatMessage(role, content string) {
	msg := models.ChatMessage{
		Role:    role,
		Content: content,
	}
	h.db.Create(&msg)
}

func (h *Handler) ListDocuments(c *gin.Context) {
	var docs []models.Document
	if err := h.db.Find(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}
	c.JSON(http.StatusOK, docs)
}

func (h *Handler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	var doc models.Document
	if err := h.db.Preload("Chunks").First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}
	c.JSON(http.StatusOK, doc)
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")

	// Get the document and its chunks before deleting
	var doc models.Document
	if err := h.db.Preload("Chunks").First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Delete from ChromaDB if available
	if h.chroma != nil && h.chroma.IsHealthy() && len(doc.Chunks) > 0 {
		var chunkIDs []uint
		for _, chunk := range doc.Chunks {
			chunkIDs = append(chunkIDs, chunk.ID)
		}
		if err := h.chroma.DeleteChunksByDocumentID(chunkIDs); err != nil {
			log.Printf("WARNING: Failed to delete chunks from ChromaDB for document %d: %v", doc.ID, err)
			// Continue with SQLite deletion even if ChromaDB delete fails
		}
	}

	// Delete from SQLite
	if err := h.db.Delete(&models.Document{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
