package services

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

type PDFProcessor struct {
	ChunkSize    int // Characters per chunk (default 800)
	ChunkOverlap int // Overlap between chunks (default 100)
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		ChunkSize:    800,
		ChunkOverlap: 100,
	}
}

type ExtractResult struct {
	Text       string   // Full text
	Pages      []string // Text per page
	TotalPages int
}

// ExtractText pulls text from PDF
func (p *PDFProcessor) ExtractText(reader io.Reader) (*ExtractResult, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	pdfReader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}

	numPages := pdfReader.NumPage()
	result := &ExtractResult{
		Pages:      make([]string, 0, numPages),
		TotalPages: numPages,
	}

	var fullText strings.Builder

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			fmt.Printf("Warning: page %d failed: %v\n", pageNum, err)
			continue
		}

		cleaned := p.cleanText(text)
		if cleaned == "" {
			continue
		}

		result.Pages = append(result.Pages, cleaned)
		fullText.WriteString(cleaned)
		fullText.WriteString("\n\n")
	}

	result.Text = fullText.String()
	return result, nil
}

func (p *PDFProcessor) cleanText(text string) string {
	// Normalize whitespace
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\t", " ")

	// Collapse multiple spaces/newlines
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	return strings.TrimSpace(text)
}

type Chunk struct {
	Content    string
	PageNumber int
	Index      int
}

// CreateChunks splits text into overlapping chunks
func (p *PDFProcessor) CreateChunks(result *ExtractResult) []Chunk {
	var chunks []Chunk
	chunkIndex := 0

	for pageNum, pageText := range result.Pages {
		pageChunks := p.chunkPage(pageText, pageNum+1, &chunkIndex)
		chunks = append(chunks, pageChunks...)
	}

	return chunks
}

func (p *PDFProcessor) chunkPage(text string, pageNum int, chunkIndex *int) []Chunk {
	var chunks []Chunk

	// Short page = single chunk
	if utf8.RuneCountInString(text) <= p.ChunkSize {
		if text != "" {
			chunks = append(chunks, Chunk{
				Content:    text,
				PageNumber: pageNum,
				Index:      *chunkIndex,
			})
			*chunkIndex++
		}
		return chunks
	}

	// Split by sentences for natural boundaries
	sentences := p.splitIntoSentences(text)
	var currentChunk strings.Builder
	currentLen := 0

	for _, sentence := range sentences {
		sentenceLen := utf8.RuneCountInString(sentence)

		// If adding this sentence exceeds limit, save chunk
		if currentLen > 0 && currentLen+sentenceLen > p.ChunkSize {
			chunks = append(chunks, Chunk{
				Content:    strings.TrimSpace(currentChunk.String()),
				PageNumber: pageNum,
				Index:      *chunkIndex,
			})
			*chunkIndex++

			// Start new chunk with overlap from previous
			overlap := p.getOverlap(currentChunk.String())
			currentChunk.Reset()
			currentChunk.WriteString(overlap)
			currentLen = utf8.RuneCountInString(overlap)
		}

		currentChunk.WriteString(sentence)
		currentLen += sentenceLen
	}

	// Don't forget last chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, Chunk{
			Content:    strings.TrimSpace(currentChunk.String()),
			PageNumber: pageNum,
			Index:      *chunkIndex,
		})
		*chunkIndex++
	}

	return chunks
}

func (p *PDFProcessor) splitIntoSentences(text string) []string {
	var sentences []string
	var current strings.Builder

	for i, r := range text {
		current.WriteRune(r)

		// Look for sentence endings
		if r == '.' || r == '!' || r == '?' {
			if i+1 >= len(text) || text[i+1] == ' ' || text[i+1] == '\n' {
				sentences = append(sentences, current.String())
				current.Reset()
			}
		}
	}

	if current.Len() > 0 {
		sentences = append(sentences, current.String())
	}

	return sentences
}

func (p *PDFProcessor) getOverlap(text string) string {
	if utf8.RuneCountInString(text) <= p.ChunkOverlap {
		return text
	}

	runes := []rune(text)
	start := len(runes) - p.ChunkOverlap
	if start < 0 {
		start = 0
	}

	return string(runes[start:])
}
