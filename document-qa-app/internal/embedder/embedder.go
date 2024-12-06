// internal/embedder/embedder.go
package embedder

import (
	"context"

	"document-qa-app/internal/loader"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
)

type Embedder struct {
	embeddings.Embedder
}

func NewOpenAIEmbedder(apiKey string) (*Embedder, error) {
	client, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
	)
	if err != nil {
		return nil, err
	}

	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, err
	}

	return &Embedder{embedder}, nil
}

func (e *Embedder) GenerateEmbeddings(documents []loader.Document) ([][]float32, error) {
	ctx := context.Background()

	// Extract text from documents
	texts := make([]string, len(documents))
	for i, doc := range documents {
		texts[i] = doc.Content
	}

	// Generate embeddings
	embeddings, err := e.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}
