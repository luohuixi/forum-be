package service

// import (
// 	"encoding/json"
// 	"forum-feed/dao"
// 	logger "forum/log"
// 	"forum/model"
// )
//
// // SubServiceRun ... 写入feed数据
// func SubServiceRun() {
// 	var feed dao.FeedModel
//
// 	ch := model.PubSubClient.Self.Channel()
// 	for msg := range ch {
// 		logger.Info("received")
//
// 		if err := json.Unmarshal([]byte(msg.Payload), &feed); err != nil {
// 			panic(err)
// 		}
//
// 		if err := feed.Create(); err != nil {
// 			logger.Error(err.Error())
// 		}
// 	}
// }
