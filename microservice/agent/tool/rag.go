package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"forum-agent/core"

	es8indexer "github.com/cloudwego/eino-ext/components/indexer/es8"
	es8retriever "github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

const (
	agentVectorIndexEnv = "AGENT_VECTOR_INDEX"
	agentESDimEnv       = "AGENT_ES_DIM"
)

type VectorStoreTool struct{}
type VectorSearchTool struct{}

type storeItem struct {
	ID      string         `json:"id"`
	Content string         `json:"content"`
	Meta    map[string]any `json:"meta,omitempty"`
}

type storeInput struct {
	Items []storeItem `json:"items"`
}

type searchInput struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

func (t *VectorStoreTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "vector_store",
		Desc: "Store one or many text documents into Elasticsearch vector database using Eino indexer. Input JSON: {\"items\":[{\"content\":\"...\",\"meta\":{\"post_id\":\"...\",\"title\":\"...\",...}}]}. meta must contain post_id and title.",
	}, nil
}

func (t *VectorStoreTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...einotool.Option) (string, error) {
	var input storeInput
	if err := json.Unmarshal([]byte(argumentsInJSON), &input); err != nil {
		return "", err
	}

	var docs []*schema.Document
	for _, item := range input.Items {
		if strings.TrimSpace(item.Content) == "" {
			continue
		}
		uid, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		item.ID = uid.String()

		doc := &schema.Document{ID: item.ID, Content: item.Content, MetaData: item.Meta}
		docs = append(docs, doc)
	}

	indexer, err := newIndexer(ctx)
	if err != nil {
		return "", err
	}
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		return "", err
	}

	out, _ := json.Marshal(map[string]any{
		"stored_ids": ids,
		"index":      vectorIndexName(),
	})
	return string(out), nil
}

func (t *VectorSearchTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "vector_search",
		Desc: "Search Elasticsearch vector database using Eino retriever. Input JSON: {\"query\":\"...\",\"top_k\":\"...\"}",
	}, nil
}

func (t *VectorSearchTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...einotool.Option) (string, error) {
	var input searchInput
	if err := json.Unmarshal([]byte(argumentsInJSON), &input); err != nil {
		return "", err
	}
	if strings.TrimSpace(input.Query) == "" {
		return "", fmt.Errorf("query is required")
	}
	if input.TopK <= 0 {
		input.TopK = 3
	}

	retriever, err := newRetriever(ctx, input.TopK)
	if err != nil {
		return "", err
	}

	docs, err := retriever.Retrieve(ctx, input.Query)
	if err != nil {
		return "", err
	}

	type item struct {
		ID      string         `json:"id"`
		Content string         `json:"content"`
		Meta    map[string]any `json:"meta,omitempty"`
	}
	items := make([]item, 0, len(docs))
	for _, doc := range docs {
		items = append(items, item{ID: doc.ID, Content: doc.Content, Meta: removeVector(doc.MetaData)})
	}

	out, _ := json.Marshal(map[string]any{
		"results":     items,
		"return_nums": len(items),
	})
	return string(out), nil
}

func newIndexer(ctx context.Context) (*es8indexer.Indexer, error) {
	return es8indexer.NewIndexer(ctx, &es8indexer.IndexerConfig{
		Client:    core.ES,
		Index:     vectorIndexName(),
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
		Embedding: core.Embedder,
		IndexSpec: &es8indexer.IndexSpec{
			Mappings: map[string]any{
				"properties": map[string]any{
					"content": map[string]any{"type": "text"},
					"content_vector": map[string]any{
						"type":       "dense_vector",
						"dims":       getDims(),
						"index":      true,
						"similarity": "cosine",
					},
				},
			},
		},
	})
}

func newRetriever(ctx context.Context, topK int) (*es8retriever.Retriever, error) {
	return es8retriever.NewRetriever(ctx, &es8retriever.RetrieverConfig{
		Client:    core.ES,
		Index:     vectorIndexName(),
		TopK:      topK,
		Embedding: core.Embedder,
		SearchMode: search_mode.SearchModeDenseVectorSimilarity(
			search_mode.DenseVectorSimilarityTypeCosineSimilarity,
			"content_vector",
		),
	})
}

func vectorIndexName() string {
	if v := os.Getenv(agentVectorIndexEnv); v != "" {
		return v
	}
	return "forum_agent_vectors"
}

func getDims() int {
	if dim := os.Getenv(agentESDimEnv); dim != "" {
		if d, err := strconv.Atoi(dim); err == nil {
			return d
		}
	}
	return 1536
}

func removeVector(origin map[string]any) map[string]any {
	meta := make(map[string]any)
	for k, v := range origin {
		if k == "content_vector" {
			continue
		}
		meta[k] = v
	}
	return meta
}

func ToolListRAG() []einotool.BaseTool {
	return []einotool.BaseTool{
		&VectorStoreTool{},
		&VectorSearchTool{},
	}
}
