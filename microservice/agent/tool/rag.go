package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"forum-agent/core"

	es8retriever "github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type VectorSearchTool struct{}

type searchInput struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

func (t *VectorSearchTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "vector_search",
		Desc: "Search Elasticsearch vector database using Eino retriever. Input JSON: {\"query\":\"...\", \"top_k\":n, top_k must be int.}",
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

func newRetriever(ctx context.Context, topK int) (*es8retriever.Retriever, error) {
	return es8retriever.NewRetriever(ctx, &es8retriever.RetrieverConfig{
		Client:    core.ES,
		Index:     core.VectorIndexName(),
		TopK:      topK,
		Embedding: core.Embedder,
		SearchMode: search_mode.SearchModeDenseVectorSimilarity(
			search_mode.DenseVectorSimilarityTypeCosineSimilarity,
			"content_vector",
		),
	})
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
		&VectorSearchTool{},
	}
}
