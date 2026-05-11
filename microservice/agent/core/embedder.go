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

// 限制batchSize最大为10，主要是防止语义分割时过多分块同时给embedder而报错
type limitedEmbedder struct {
	embedding.Embedder
	maxBatchSize int
}

func (e *limitedEmbedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	if len(texts) <= e.maxBatchSize {
		return e.Embedder.EmbedStrings(ctx, texts, opts...)
	}

	var vectors [][]float64
	for start := 0; start < len(texts); start += e.maxBatchSize {
		end := start + e.maxBatchSize
		if end > len(texts) {
			end = len(texts)
		}

		batchVectors, err := e.Embedder.EmbedStrings(ctx, texts[start:end], opts...)
		if err != nil {
			return nil, err
		}
		vectors = append(vectors, batchVectors...)
	}

	return vectors, nil
}

func NewLimitedEmbedder(embedder embedding.Embedder, maxBatchSize int) embedding.Embedder {
	if maxBatchSize <= 0 {
		maxBatchSize = 10
	}
	return &limitedEmbedder{Embedder: embedder, maxBatchSize: maxBatchSize}
}
