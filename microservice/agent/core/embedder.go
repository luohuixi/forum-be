package core

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/spf13/viper"
)

var Embedder embedding.Embedder

func EmbedderInit() {
	ctx := context.Background()
	var err error
	dim := viper.GetInt("embed.dimensions")
	var dimensions *int
	if dim > 0 {
		dimensions = &dim
	}

	Embedder, err = dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		APIKey:     viper.GetString("embed.api_key"),
		Model:      viper.GetString("embed.model_name"),
		Dimensions: dimensions,
	})
	if err != nil {
		log.Fatalf("embedder.NewEmbedder: %v", err)
	}
}
