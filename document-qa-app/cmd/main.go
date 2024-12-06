// Project Structure:
// /document-qa-app
// ├── cmd/
// │   └── main.go
// ├── internal/
// │   ├── loader/
// │   │   └── loader.go
// │   ├── embedder/
// │   │   └── embedder.go
// │   ├── vectorstore/
// │   │   └── chroma.go
// │   └── chatbot/
// │       └── chatbot.go
// ├── go.mod
// └── README.md

// cmd/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"document-qa-app/internal/chatbot"
	"document-qa-app/internal/loader"
	"document-qa-app/internal/vectorstore"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

type DataSource struct {
	Type   DataSourceType
	Source string
}
type DataSourceType int

const (
	SOURCE_WEBSITE DataSourceType = iota
	SOURCE_JSON_FILE
)

func main() {
	// Configuration and dependency injection
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dataSources := []DataSource{
		{Type: SOURCE_WEBSITE, Source: "https://genk.vn/bo-ra-30-trieu-mua-yamaha-nvx-cu-nhung-van-phai-thanh-toan-them-hon-100-trieu-vi-nhung-dieu-nay-2024120415032637.chn"},
	}

	config := struct {
		OpenAIAPIKey     string
		ChromaHost       string
		QdrantUrl        string
		QdrantAPIKey     string
		QdrantCollection string
		DataSource       []DataSource
	}{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		ChromaHost:   os.Getenv("CHROMA_URL"),
		// Qdrant
		QdrantUrl:        os.Getenv("QDRANT_URL"),
		QdrantAPIKey:     os.Getenv("QDRANT_API_KEY"),
		QdrantCollection: os.Getenv("QDRANT_COLLECTION"),
		// DataSource:
		DataSource: dataSources,
	}

	// Initialize OpenAI LLM
	// llm, err := openai.New(
	// 	openai.WithModel("gpt-3.5-turbo"),
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// Initialize Ollama LLM
	ollamaLLM, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		log.Fatal(err)
	}
	ollamaEmbedder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		log.Fatal(err)
	}

	// Load document
	var chunks []loader.Document
	for _, dataSource := range config.DataSource {
		switch dataSource.Type {
		case SOURCE_WEBSITE:
			docs, err := loader.LoadFromURL(dataSource.Source)
			if err != nil {
				log.Fatalf("Failed to load document from URL: %v", err)
			}
			// Split documents into chunks
			chunks = append(chunks, loader.SplitDocuments(docs)...)
		}
	}

	// // Initialize embedder
	// embedder, err := embedder.NewOpenAIEmbedder(config.OpenAIAPIKey)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize embedder: %v", err)
	// }

	// // Generate embeddings
	// embeddings, err := embedder.GenerateEmbeddings(chunks)
	// if err != nil {
	// 	log.Fatalf("Failed to generate embeddings: %v", err)
	// }

	// Initialize vector store
	vectorStore, err := vectorstore.NewQdrantStore(config.QdrantUrl, config.QdrantAPIKey, config.QdrantCollection, ollamaEmbedder)
	if err != nil {
		log.Fatalf("Failed to initialize vector store: %v", err)
	}

	// Store embeddings in vector database
	err = vectorStore.AddDocuments(chunks)
	if err != nil {
		log.Fatalf("Failed to add documents to vector store: %v", err)
	}

	// Initialize chatbot
	bot, err := chatbot.NewChatbot(ollamaLLM, vectorStore, context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize chatbot: %v", err)
	}

	// Console UI for Q&A
	consoleUI(bot)

}

func consoleUI(bot *chatbot.Chatbot) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter your question (or 'exit' to quit): ")
		if scanner.Scan() {
			query := scanner.Text()
			if query == "exit" {
				break
			}
			// Retrieve and answer question
			answer, err := bot.Answer(query)
			if err != nil {
				fmt.Printf("Error answering question: %v\n", err)
				continue
			}
			fmt.Println("Answer:", answer)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
			break
		}
	}
}
