// internal/vectorstore/chroma.go
package vectorstore

import (
	"context"
	"log"
	"net/url"

	"document-qa-app/internal/loader"

	qdrlient "github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

type QdrantStore struct {
	store *qdrant.Store
}

func NewQdrantStore(qdrantUrl string, qdrantAPIKey string, qdrantCollectionName string, e embeddings.Embedder) (*QdrantStore, error) {
	client, err := qdrlient.NewClient(&qdrlient.Config{
		Host:   "e471ec23-529e-48d1-84bc-8abd79d0c03b.us-west-2-0.aws.cloud.qdrant.io",
		Port:   6334,
		APIKey: qdrantAPIKey,
		UseTLS: true,
		// uses default config with minimum TLS version set to 1.3
		// TLSConfig: &tls.Config{...},
		// GrpcOptions: []grpc.DialOption{},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	if existed, err := client.CollectionExists(context.Background(), qdrantCollectionName); !existed {
		err := client.CreateCollection(context.Background(), &qdrlient.CreateCollection{
			CollectionName: qdrantCollectionName,
			VectorsConfig: qdrlient.NewVectorsConfig(&qdrlient.VectorParams{
				Size:     3072,
				Distance: qdrlient.Distance_Cosine,
			}),
		})
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	// Mew vector store
	url, err := url.Parse(qdrantUrl)
	if err != nil {
		log.Fatal(err)
	}
	store, err := qdrant.New(
		qdrant.WithURL(*url),
		qdrant.WithAPIKey(qdrantAPIKey),
		qdrant.WithCollectionName(qdrantCollectionName),
		qdrant.WithEmbedder(e),
	)
	if err != nil {
		return nil, err
	}

	return &QdrantStore{store: &store}, nil
}

func (c *QdrantStore) AddDocuments(documents []loader.Document) error {
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

func (c *QdrantStore) Search(ctx context.Context, query string, k int) ([]loader.Document, error) {
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
