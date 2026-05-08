package service

import (
	"context"
	"encoding/json"
	"fmt"
	"forum-agent/core"
	pb "forum-agent/proto"
	"forum/log"
	"forum/pkg/agentctx"
	"forum/pkg/tracer"

	"go.uber.org/zap"
)

func (a *AgentService) GiveAnswer(ctx context.Context, req *pb.GiveAnswerRequest, _ *pb.Response) error {
	log.Info(fmt.Sprintf("Receive request to give answer for post(%v)", req.PostId), zap.String("trace_id", tracer.GetTraceId(ctx)))

	// 不等待
	go func() {
		baseCtx := context.WithoutCancel(ctx)
		newCtx := agentctx.WithTokenUsage(baseCtx)
		content, err := a.formatRequest(req)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to format request for post(%v)", req.PostId), zap.String("trace_id", tracer.GetTraceId(newCtx)), zap.Error(err))
			return
		}

		prompt, err := core.CommentPrompt(string(content))
		if err != nil {
			log.Error(fmt.Sprintf("Failed to create prompt for post(%v)", req.PostId), zap.String("trace_id", tracer.GetTraceId(newCtx)), zap.Error(err))
			return
		}

		_, err = a.agent.Generate(newCtx, prompt)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to give answer for post(%v)", req.PostId), zap.String("trace_id", tracer.GetTraceId(newCtx)), zap.Error(err))
			log.Info("agent token summary", zap.Any("summary", agentctx.FinalTokenLogFields(newCtx)), zap.String("trace_id", tracer.GetTraceId(newCtx)))
			return
		}

		log.Info("agent token summary", zap.Any("summary", agentctx.FinalTokenLogFields(newCtx)), zap.String("trace_id", tracer.GetTraceId(newCtx)))
	}()

	return nil
}

func (a *AgentService) formatRequest(req *pb.GiveAnswerRequest) ([]byte, error) {
	post, err := a.Dao.GetPostById(req.PostId)
	if err != nil {
		return nil, err
	}

	var format FormatRequest
	format.Title = post.Title
	format.Content = post.Content
	format.PostId = post.Id
	format.ExtraContent = req.ExtraContent

	return json.Marshal(format)
}
