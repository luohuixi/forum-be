package service

import (
	"context"
	"encoding/json"
	"forum-feed/dao"
	logger "forum/log"
	"forum/model"
)

// SubServiceRun ... 写入feed数据
func SubServiceRun() {
	var feed dao.FeedModel

	ctx := context.Background()
	for {
		m, err := model.KafkaReader.FetchMessage(ctx)
		if err != nil {
			logger.Error(err.Error())
			break
		}
		if err := json.Unmarshal(m.Value, &feed); err != nil {
			panic(err)
		}

		logger.Info(string(m.Value))
		if err := feed.Create(); err != nil {
			logger.Error(err.Error())
		}

		if err := model.KafkaReader.CommitMessage(ctx, m); err != nil {
			logger.Error(err.Error())
			break
		}
	}

}
