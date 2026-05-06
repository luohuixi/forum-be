package service

import (
	"context"
	"encoding/json"
	"fmt"
	"forum-agent/core"
	pb "forum-agent/proto"

	"forum/log"

	"go.uber.org/zap"
)

func (a *AgentService) AddKnowledge(ctx context.Context, req *pb.AddKnowledgeRequest, resp *pb.Response) error {
	log.Info("Receive request to add knowledge", zap.Any("knowledge", req.PostId))

	post, err := a.Dao.GetPostById(req.PostId)
	if err != nil {
		return err
	}
	postJson, err := json.Marshal(post)
	if err != nil {
		return err
	}

	prompt, err := core.StorePrompt(string(postJson))
	if err != nil {
		return err
	}
	ans, err := a.agent.Generate(ctx, prompt)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Successfully store post(%v) into es", req.PostId), zap.Any("final_reply", ans.Content))
	return nil
}
