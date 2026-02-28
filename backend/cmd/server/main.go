package main

import (
	"log"

	"azeru/internal/config"
	"azeru/internal/database"
	"azeru/internal/handlers"
	"azeru/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	log.Println("Starting Enterprise Brain...")
	log.Printf("Environment: %s", cfg.Environment)

	// Connect to database
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Initialize Ollama client (optional - for embeddings)
	ollamaClient := services.NewOllamaClient(cfg.OllamaURL)
	if ollamaClient == nil {
		log.Println("WARNING: Ollama not available - running without embeddings")
	} else if !ollamaClient.IsHealthy() {
		log.Println("WARNING: Ollama not healthy - running without embeddings")
		ollamaClient = nil
	} else {
		log.Println("Ollama connected successfully")
	}

	// Initialize ChromaDB client (optional - for vector search)
	chromaClient := services.NewChromaClient(cfg.ChromaURL)
	if chromaClient == nil {
		log.Println("WARNING: ChromaDB not available - falling back to FTS5 only")
	} else if !chromaClient.IsHealthy() {
		log.Println("WARNING: ChromaDB not healthy - falling back to FTS5 only")
		chromaClient = nil
	} else {
		log.Println("ChromaDB connected successfully")
	}

	// Create handler with all dependencies
	h := handlers.NewHandler(cfg, db, ollamaClient, chromaClient)

	// Setup router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := router.Group("/api")
	{
		api.POST("/upload", h.UploadDocument)
		api.GET("/documents", h.ListDocuments)
		api.GET("/documents/:id", h.GetDocument)
		api.DELETE("/documents/:id", h.DeleteDocument)
		api.POST("/chat", h.Chat)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("Server running on http://%s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
