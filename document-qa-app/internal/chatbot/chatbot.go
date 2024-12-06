// internal/chatbot/chatbot.go
package chatbot

import (
	"context"

	"document-qa-app/internal/vectorstore"

	"github.com/tmc/langchaingo/llms"
)

type Chatbot struct {
	llm         llms.Model
	vectorStore vectorstore.VectorStore
	// embedder    embeddings.Embedder
	ctx context.Context
}

func NewChatbot(
	llm llms.Model,
	vectorStore vectorstore.VectorStore,
	// embedder *embedder.Embedder,
	ctx context.Context,
) (*Chatbot, error) {
	return &Chatbot{
		vectorStore: vectorStore,
		// embedder:    embedder,
		llm: llm,
		ctx: ctx,
	}, nil
}

func (c *Chatbot) Answer(query string) (string, error) {
	ctx := context.Background()
	// Retrieve relevant documents
	relevantDocs, err := c.vectorStore.Search(c.ctx, query, 3)
	if err != nil {
		return "", err
	}

	// Prepare context for LLM
	var context string
	for _, doc := range relevantDocs {
		context += doc.Content + "\n\n"
	}

	// Construct prompt
	prompt := "Context: " + context + "\n\nQuestion: " + query + "\n\nAnswer:"

	// Generate answer using LLM
	response, err := c.llm.Call(ctx, prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
