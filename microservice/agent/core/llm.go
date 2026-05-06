package core

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/spf13/viper"
)

var ChatModel *openai.ChatModel

func ChatModelInit() {
	ctx := context.Background()
	var err error

	ChatModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  viper.GetString("llm.api_key"),
		Model:   viper.GetString("llm.model_id"),
		BaseURL: viper.GetString("llm.base_url"),
	})
	if err != nil {
		log.Fatalf("NewChatModel failed, err=%v", err)
	}
}
