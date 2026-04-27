package main

import (
	"context"
	"forum-agent/async"
	"forum-agent/core"
	"forum-agent/dao"
	"log"

	"forum-agent/tool"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	_ = godotenv.Load()
	viper.SetDefault("kafka.username", "")
	viper.SetDefault("kafka.password", "")
	viper.SetDefault("kafka.addr", "localhost:9092")
	viper.SetDefault("db.username", "forum")
	viper.SetDefault("db.password", "Muxistudio517")
	viper.SetDefault("db.addr", "rm-bp1pkcz9y269h30fm5o.mysql.rds.aliyuncs.com")
	viper.SetDefault("db.name", "forum")
	viper.SetDefault("kafka.topic", "testA")
	viper.SetDefault("kafka.group_id", "group1")
}

func main() {
	core.ChatModelInit()
	core.EmbedderInit()
	core.ESClientInit()
	//	core.KafkaInit()
	dao.Init()
	manager, err := async.NewAsyncManager(dao.GetDao())
	if err != nil {
		log.Fatal(err)
	}
	manager.Run()

	select {}
}
