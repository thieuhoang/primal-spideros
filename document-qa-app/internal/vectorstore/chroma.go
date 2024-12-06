// internal/vectorstore/chroma.go
package vectorstore

import (
	"context"
	"os"

	"document-qa-app/internal/loader"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type ChromaStore struct {
	store *chroma.Store
}

func NewChromaStore(host string, e embeddings.Embedder) (*ChromaStore, error) {
	store, err := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithEmbedder(e),
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace(uuid.New().String()),
		// chroma.WithCollectionName("document_collection"),
	)
	if err != nil {
		return nil, err
	}

	return &ChromaStore{store: &store}, nil
}

func (c *ChromaStore) AddDocuments(documents []loader.Document) error {
	ctx := context.Background()

	// Prepare documents for vector store
	docs := make([]schema.Document, len(documents))
	for i, doc := range documents {
		docs[i] = schema.Document{
			PageContent: doc.Content,
			Metadata:    doc.Metadata,
		}
	}

	// Add documents with embeddings
	_, err := c.store.AddDocuments(ctx, docs)
	return err
}

func (c *ChromaStore) Search(ctx context.Context, query string, k int) ([]loader.Document, error) {
	// Perform similarity search
	results, err := c.store.SimilaritySearch(ctx, query, k)
	if err != nil {
		return nil, err
	}

	// Convert results to our Document type
	documents := make([]loader.Document, len(results))
	for i, result := range results {
		documents[i] = loader.Document{
			Content:  result.PageContent,
			Metadata: result.Metadata,
		}
	}

	return documents, nil
}
