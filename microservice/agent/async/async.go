package async

import (
	"context"
	"forum-agent/core"
	"forum-agent/dao"
	"time"
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

func NewAsyncManager(dao dao.Interface, agent *core.ReActAgent) *AsyncManager {
	return &AsyncManager{
		dao:   dao,
		agent: agent,
	}
}

func (a *AsyncManager) Begin() {
	ctx := context.Background()
	go a.CreateCommentFromKafka(ctx)
}
