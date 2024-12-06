// internal/vectorstore/chroma.go
package vectorstore

import (
	"context"
	"document-qa-app/internal/loader"
)

type VectorStore interface {
	AddDocuments(documents []loader.Document) error
	Search(ctx context.Context, query string, k int) ([]loader.Document, error)
}
