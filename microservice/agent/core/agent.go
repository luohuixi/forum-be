package core

import (
	"context"
	"fmt"
	"log"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

const MaxStep = 8

type ReActAgent struct {
	Inner *react.Agent
}

var reActAgent *ReActAgent

func GetReActAgent() *ReActAgent {
	return reActAgent
}

func NewReActAgent(ctx context.Context, tools []einotool.BaseTool) (*ReActAgent, error) {
	if ChatModel == nil {
		return nil, fmt.Errorf("chat model is not initialized")
	}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: ChatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
		MaxStep: MaxStep,
	})
	if err != nil {
		return nil, err
	}

	return &ReActAgent{Inner: agent}, nil
}

func (a *ReActAgent) Generate(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	if a == nil || a.Inner == nil {
		return nil, fmt.Errorf("agent is nil")
	}
	return a.Inner.Generate(ctx, messages)
}

func AgentInit(tools []einotool.BaseTool) {
	var err error
	reActAgent, err = NewReActAgent(context.Background(), tools)
	if err != nil {
		log.Fatal(err)
	}
}
