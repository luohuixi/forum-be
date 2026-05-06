package async

import (
	"context"
	"encoding/json"
	"fmt"
	"forum-agent/core"
	"forum-agent/dao"
	"forum/log"
	"time"

	"go.uber.org/zap"
)

func (a *AsyncManager) CreateCommentFromKafka(ctx context.Context) {
	reader := core.KafkaReader()
	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Error("Failed to receive message from Kafka", zap.Error(err))
			return
		}
		commentJson := msg.Value

		var c dao.CommentAgentReturn
		if err := json.Unmarshal(commentJson, &c); err != nil {
			log.Error("Failed to unmarshal comment JSON", zap.Error(err))
			return
		}

		cm := dao.ChangeToCommentModel(&c)
		for i := 0; i < MaxTries; i++ {
			err := cm.Create()
			if err != nil {
				log.Error(fmt.Sprintf("Failed to create comment(%v)", string(commentJson)), zap.Error(err))
				time.Sleep(WaitForError)
			} else {
				fmt.Println("Create comment successfully", c)
				break
			}
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Error(fmt.Sprintf("Failed to commit messages(%v)", string(commentJson)), zap.Error(err))
		}
	}
}
