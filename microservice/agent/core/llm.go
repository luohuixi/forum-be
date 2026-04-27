package core

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

var ChatModel *openai.ChatModel

func ChatModelInit() {
	ctx := context.Background()
	var err error

	ChatModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("LLM_API_KEY"),
		Model:   os.Getenv("LLM_MODEL_ID"),
		BaseURL: os.Getenv("LLM_BASE_URL"),
	})
	if err != nil {
		log.Fatalf("NewChatModel failed, err=%v", err)
	}
}
