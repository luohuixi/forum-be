package service

import (
	"context"
	"encoding/json"
	"forum-agent/core"
	pb "forum-agent/proto"

	"forum/log"
	"forum/pkg/agentctx"
	"forum/pkg/tracer"

	"go.uber.org/zap"
)

func (a *AgentService) AddKnowledge(ctx context.Context, req *pb.AddKnowledgeRequest, resp *pb.Response) error {
	log.Info("Receive request to add knowledge", zap.Any("knowledge", req.PostId))
	ctx = agentctx.WithTokenUsage(ctx)

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

	_, err = a.agent.Generate(ctx, prompt)
	if err != nil {
		log.Info("agent token summary", zap.Any("summary", agentctx.FinalTokenLogFields(ctx)), zap.String("trace_id", tracer.GetTraceId(ctx)))
		return err
	}

	log.Info("agent token summary", zap.Any("summary", agentctx.FinalTokenLogFields(ctx)), zap.String("trace_id", tracer.GetTraceId(ctx)))
	return nil
}
