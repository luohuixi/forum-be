package service

import (
	"forum-agent/core"
	"forum-agent/dao"
)

// AgentService ... 智能体服务
type AgentService struct {
	Dao   dao.Interface
	agent *core.ReActAgent
}

func New(i dao.Interface, agent *core.ReActAgent) *AgentService {
	service := new(AgentService)
	service.Dao = i
	service.agent = agent

	return service
}

type FormatRequest struct {
	PostId       uint32 `json:"post_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ExtraContent string `json:"extra_content"`
}
