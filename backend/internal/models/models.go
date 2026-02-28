package models

import (
	"time"

	"gorm.io/gorm"
)

// Document = one uploaded PDF
type Document struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Status      string `json:"status"` // processing, completed, failed

	// One Document has many Chunks
	Chunks     []Chunk `json:"chunks,omitempty" gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE;"`
	ChunkCount int     `json:"chunk_count"`
}

// Chunk = searchable piece of text from PDF
type Chunk struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	DocumentID uint   `json:"document_id"`
	Content    string `gorm:"type:text" json:"content"`

	// Temporary: Store embedding as JSON until Phase 2 (ChromaDB)
	Embedding []byte `gorm:"type:blob" json:"-"`

	ChunkIndex int `json:"chunk_index"`
	PageNum    int `json:"page_num"`
}

// ChatMessage = conversation history
type ChatMessage struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Role         string    `json:"role"` // "user" or "assistant"
	Content      string    `json:"content"`
	SourceChunks string    `json:"source_chunks,omitempty" gorm:"type:text"` // JSON array
}
