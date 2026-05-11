package core

import (
	"context"

	es8indexer "github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/schema"
	"github.com/spf13/viper"
)

func NewIndexer(ctx context.Context) (*es8indexer.Indexer, error) {
	return es8indexer.NewIndexer(ctx, &es8indexer.IndexerConfig{
		Client:    ES,
		Index:     VectorIndexName(),
		BatchSize: 5,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (map[string]es8indexer.FieldValue, error) {
			fields := map[string]es8indexer.FieldValue{
				"content": {Value: doc.Content, EmbedKey: "content_vector"},
			}
			for k, v := range doc.MetaData {
				fields[k] = es8indexer.FieldValue{Value: v}
			}
			return fields, nil
		},
		Embedding: Embedder,
		IndexSpec: &es8indexer.IndexSpec{
			Mappings: map[string]any{
				"properties": map[string]any{
					"content": map[string]any{"type": "text"},
					"content_vector": map[string]any{
						"type":       "dense_vector",
						"dims":       VectorDims(),
						"index":      true,
						"similarity": "cosine",
					},
				},
			},
		},
	})
}

func VectorIndexName() string {
	if v := viper.GetString("agent_es.index"); v != "" {
		return v
	}
	return "forum_agent_vectors"
}

func VectorDims() int {
	if dim := viper.GetInt("agent_es.dim"); dim > 0 {
		return dim
	}
	if dim := viper.GetInt("embed.dimensions"); dim > 0 {
		return dim
	}
	return 1024
}
