// pdf_service.go - This is the wrapper for the handler
package services

import (
	"fmt"
	"os"
)

// PDFService wraps PDFProcessor to provide a simple interface for handlers
// Handler needs: ExtractText(filepath) (string, error)
//
//	ChunkText(filepath, size, overlap) []string
type PDFService struct {
	processor *PDFProcessor
}

func NewPDFService() *PDFService {
	return &PDFService{
		processor: NewPDFProcessor(),
	}
}

// ExtractText opens a file and returns plain text
func (s *PDFService) ExtractText(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer file.Close()

	result, err := s.processor.ExtractText(file)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}

// ChunkText opens a file and returns chunk contents as strings
func (s *PDFService) ChunkText(filePath string, chunkSize, overlap int) []string {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Temporarily override processor settings
	oldChunkSize := s.processor.ChunkSize
	oldOverlap := s.processor.ChunkOverlap
	s.processor.ChunkSize = chunkSize
	s.processor.ChunkOverlap = overlap
	defer func() {
		s.processor.ChunkSize = oldChunkSize
		s.processor.ChunkOverlap = oldOverlap
	}()

	result, err := s.processor.ExtractText(file)
	if err != nil {
		return nil
	}

	chunks := s.processor.CreateChunks(result)

	// Extract just the content strings
	var contents []string
	for _, chunk := range chunks {
		contents = append(contents, chunk.Content)
	}
	return contents
}
