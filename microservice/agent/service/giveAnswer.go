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

func (a *AgentService) GiveAnswer(ctx context.Context, req *pb.GiveAnswerRequest, _ *pb.Response) error {
	log.Info(fmt.Sprintf("Receive request to give answer for post(%v)", req.PostId))

	// 不等待
	go func() {
		content, err := a.formatRequest(req)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to format request for post(%v)", req.PostId), zap.Error(err))
			return
		}

		ans, err := a.run(string(content))
		if err != nil {
			log.Error(fmt.Sprintf("Failed to give answer for post(%v)", req.PostId), zap.Error(err))
			return
		}

		log.Info(fmt.Sprintf("Successfully generate answer for post(%v)", req.PostId), zap.Any("final_reply", ans))
	}()

	return nil
}

func (a *AgentService) run(content string) (string, error) {
	ctx := context.Background()

	prompt, err := core.CommentPrompt(content)
	if err != nil {
		return "", err
	}

	ans, err := a.agent.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	return ans.Content, nil
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
