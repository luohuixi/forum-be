package client

import (
	pbc "forum-chat/proto"
	"forum/pkg/identity"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	micro "go-micro.dev/v4"
)

var ChatClient pbc.ChatService

func ChatInit(service micro.Service) {
	ChatClient = pbc.NewChatService(identity.Prefix()+"forum.service.chat", service.Client())
}
