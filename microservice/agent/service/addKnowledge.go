package service

import (
	"context"
	"fmt"

	"forum-agent/core"
	pb "forum-agent/proto"

	"forum/log"

	"github.com/cloudwego/eino/schema"
)

func (a *AgentService) AddKnowledge(ctx context.Context, req *pb.AddKnowledgeRequest, resp *pb.Response) error {
	log.Info("Receive request to add knowledge")

	var docs []*schema.Document
	var err error
	if req.SplitType == "markdown" {
		docs, err = core.MarkdownSplitter(ctx, req.PostId, req.Content, req.SplitSize)
	} else {
		docs, err = core.SemanticSpliter(ctx, req.PostId, req.Content, req.SplitSize)
	}
	if err != nil {
		return err
	}
	if len(docs) == 0 {
		return fmt.Errorf("empty markdown content")
	}

	indexer, err := core.NewIndexer(ctx)
	if err != nil {
		return err
	}
	_, err = indexer.Store(ctx, docs)
	return err
}
