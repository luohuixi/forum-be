package core

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino/components/embedding"
)

var Embedder embedding.Embedder

func EmbedderInit() {
	ctx := context.Background()
	var err error

	Embedder, err = dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		APIKey: os.Getenv("EMBED_API_KEY"),
		Model:  os.Getenv("EMBED_MODEL_NAME"),
	})
	if err != nil {
		log.Fatalf("embedder.NewEmbedder: %v", err)
	}
}
