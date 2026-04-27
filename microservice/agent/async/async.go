package async

import (
	"context"
	"forum-agent/core"
	"forum-agent/dao"
	"forum-agent/tool"
	"forum/log"
	"time"

	"go.uber.org/zap"
)

const (
	WaitForError = time.Second * 5
	MaxTries     = 3
)

type AsyncManager struct {
	dao   dao.Interface
	agent *core.ReActAgent
}

type Async interface {
	CreateCommentFromKafka(ctx context.Context)
}

func NewAsyncManager(dao dao.Interface) (*AsyncManager, error) {
	agent, err := core.NewReActAgent(context.Background(), tool.ToolList())
	if err != nil {
		log.Error("Failed to create react agent", zap.Error(err))
		return nil, err
	}
	return &AsyncManager{
		dao:   dao,
		agent: agent,
	}, nil
}

func (a *AsyncManager) Begin(ctx context.Context) {
	go a.CreateCommentFromKafka(ctx)
}
