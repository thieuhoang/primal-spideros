// internal/loader/loader.go
package loader

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

type Document struct {
	Content  string
	Metadata map[string]any
}

func LoadFromURL(url string) ([]Document, error) {
	// Download document from URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read document content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Detect document type and use appropriate loader
	// For simplicity, we'll start with basic text loading
	loader := documentloaders.NewText(strings.NewReader(string(body)))
	docs, err := loader.Load(context.Background())
	if err != nil {
		return nil, err
	}

	// Convert to our Document type
	var result []Document
	for _, doc := range docs {
		result = append(result, Document{
			Content:  doc.PageContent,
			Metadata: doc.Metadata,
		})
	}

	return result, nil
}

func SplitDocuments(docs []Document) []Document {
	// Use recursive character text splitter
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(10000),
		textsplitter.WithChunkOverlap(2000),
	)

	var chunks []Document
	for _, doc := range docs {
		splits, err := splitter.SplitText(doc.Content)
		if err != nil {
			continue
		}

		for _, split := range splits {
			chunks = append(chunks, Document{
				Content:  split,
				Metadata: doc.Metadata,
			})
		}
	}

	return chunks
}
