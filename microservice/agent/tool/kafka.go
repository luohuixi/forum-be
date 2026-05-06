package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"forum-agent/core"
	"forum-agent/dao"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type CommentQueueTool struct{}

type commentQueueInput struct {
	Content string `json:"content"`
	PostId  uint32 `json:"post_id"`
}

func (t *CommentQueueTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "comment_queue",
		Desc: "Send agent-generated comment to Kafka queue. Input JSON: {\"content\":\"...\",\"post_id\":n, post_id must be uint32.}",
	}, nil
}

func (t *CommentQueueTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...einotool.Option) (string, error) {
	var input commentQueueInput
	fmt.Println(argumentsInJSON)
	if err := json.Unmarshal([]byte(argumentsInJSON), &input); err != nil {
		return "", err
	}
	if strings.TrimSpace(input.Content) == "" {
		return "", fmt.Errorf("content is required")
	}
	if input.PostId == 0 {
		return "", fmt.Errorf("post_id is required")
	}

	payload := dao.CommentAgentReturn{
		Content: input.Content,
		PostId:  input.PostId,
	}
	msg, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	if err := core.KafkaPublish(msg); err != nil {
		return "", err
	}

	return string(msg), nil
}

func ToolListKafka() []einotool.BaseTool {
	return []einotool.BaseTool{
		&CommentQueueTool{},
	}
}
