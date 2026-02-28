package database

import (
	"fmt"
	"log"
	"strings"

	"azeru/internal/config"
	"azeru/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init initializes database connection and returns it
func Init(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Info
	if cfg.Environment == "production" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	log.Println("Connected to SQLite:", cfg.DatabasePath)

	// Set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1) // SQLite prefers single connection

	// Create tables
	if err := db.AutoMigrate(&models.Document{}, &models.Chunk{}, &models.ChatMessage{}); err != nil {
		return nil, err
	}

	// Setup full-text search
	if err := setupFTS5(db); err != nil {
		return nil, err
	}

	return db, nil
}

func setupFTS5(db *gorm.DB) error {
	// Check if FTS5 exists
	var count int64
	db.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='chunks_fts'").Scan(&count)
	if count > 0 {
		return nil
	}

	log.Println("Setting up FTS5 search index...")

	// Create virtual table for fast text search
	db.Exec(`CREATE VIRTUAL TABLE chunks_fts USING fts5(
		content,
		content='chunks',
		content_rowid='id'
	)`)

	// Triggers keep search index in sync
	db.Exec(`CREATE TRIGGER chunks_ai AFTER INSERT ON chunks BEGIN
		INSERT INTO chunks_fts(rowid, content) VALUES (new.id, new.content);
	END`)

	db.Exec(`CREATE TRIGGER chunks_ad AFTER DELETE ON chunks BEGIN
		INSERT INTO chunks_fts(chunks_fts, rowid, content) VALUES ('delete', old.id, old.content);
	END`)

	db.Exec(`CREATE TRIGGER chunks_au AFTER UPDATE ON chunks BEGIN
		INSERT INTO chunks_fts(chunks_fts, rowid, content) VALUES ('delete', old.id, old.content);
		INSERT INTO chunks_fts(rowid, content) VALUES (new.id, new.content);
	END`)

	log.Println("FTS5 ready")
	return nil
}

// sanitizeFTS5Query cleans user input so it is safe for FTS5 MATCH.
func sanitizeFTS5Query(query string) string {
	// Remove FTS5 special characters
	replacer := strings.NewReplacer(
		"\"", "",
		"*", "",
		"?", "",
		"(", "",
		")", "",
		"+", "",
		"-", "",
		"^", "",
		":", "",
		"{", "",
		"}", "",
		"~", "",
	)
	cleaned := replacer.Replace(query)

	// Split into words and quote each one
	words := strings.Fields(cleaned)
	var quoted []string
	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) > 0 {
			quoted = append(quoted, "\""+w+"\"")
		}
	}

	return strings.Join(quoted, " ")
}

// SearchChunks finds relevant chunks using full-text search
func SearchChunks(db *gorm.DB, query string, limit int) ([]models.Chunk, error) {
	var chunks []models.Chunk

	sanitized := sanitizeFTS5Query(query)

	// If sanitization removed all terms, fall back to LIKE search
	if sanitized == "" {
		sql := `SELECT * FROM chunks WHERE content LIKE ? LIMIT ?`
		err := db.Raw(sql, "%"+query+"%", limit).Scan(&chunks).Error
		return chunks, err
	}

	sql := `SELECT c.* FROM chunks c
		JOIN chunks_fts fts ON c.id = fts.rowid
		WHERE chunks_fts MATCH ?
		ORDER BY bm25(chunks_fts)
		LIMIT ?`

	err := db.Raw(sql, sanitized, limit).Scan(&chunks).Error
	return chunks, err
}
