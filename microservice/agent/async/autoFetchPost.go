package async

import (
	"context"
	"encoding/json"
	"forum-agent/core"
	"forum-agent/dao"
	"forum/log"
	"forum/util"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func (a *AsyncManager) AutoFetchPost() {
	c := cron.New()
	_, _ = c.AddFunc("0 2 * * *", func() {
		t := util.GetTimeDay(-1)
		post := a.fetch(t)
		a.Run(post)
	})
	c.Start()
}

func (a *AsyncManager) fetch(t string) string {
	var err error
	var posts *[]dao.PostModel

	for i := 0; i < MaxTries; i++ {
		posts, err = a.dao.GetValuablePost()
		if err != nil {
			log.Error("Failed to fetch posts", zap.Error(err))
			time.Sleep(WaitForError)
		} else {
			break
		}
	}

	postJson, err := json.Marshal(posts)
	if err != nil {
		log.Error("Failed to marshal posts", zap.Error(err))
		return ""
	}

	return string(postJson)
}

func (a *AsyncManager) Run(postJson string) {
	ctx := context.Background()
	prompt, err := core.StorePrompt(postJson)
	if err != nil {
		log.Error("Failed to store post", zap.Error(err))
		return
	}
	ans, err := a.agent.Generate(ctx, prompt)
	if err != nil {
		log.Error("Failed to generate post", zap.Error(err))
	} else {
		log.Info("Successfully store post into es", zap.Any("final_reply", ans.Content))
	}
}
