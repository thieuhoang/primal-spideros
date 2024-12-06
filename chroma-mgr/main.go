package main

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

type ChromaStore struct {
	store *chroma.Store
	embdr embeddings.Embedder
}

func NewChromaStore(host string) (*ChromaStore, error) {
	ollamaLLM, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		log.Fatal(err)
	}
	ollamaEmbedder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		log.Fatal(err)
	}

	store, err := chroma.New(
		chroma.WithChromaURL(os.Getenv("CHROMA_URL")),
		chroma.WithEmbedder(ollamaEmbedder),
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace(uuid.New().String()),
	)
	if err != nil {
		return nil, err
	}

	return &ChromaStore{store: &store, embdr: ollamaEmbedder}, nil
}

func (c *ChromaStore) Search(ctx context.Context, query string, k int) ([]schema.Document, error) {
	return c.store.SimilaritySearch(ctx, query, k)
}

func main() {
	_, err := NewChromaStore("http://localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

}
