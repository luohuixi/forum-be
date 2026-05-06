package client

import (
	pba "forum-agent/proto"
	"forum/pkg/identity"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	micro "go-micro.dev/v4"
)

var AgentClient pba.AgentService

func AgentInit(service micro.Service) {
	AgentClient = pba.NewAgentService(identity.Prefix()+"forum.service.agent", service.Client())
}
