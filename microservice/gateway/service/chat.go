package service

import (
	pbc "forum-chat/proto"
	"forum/pkg/handler"

	_ "github.com/micro/go-plugins/registry/kubernetes"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var ChatService micro.Service
var ChatClient pbc.ChatServiceClient

func ChatInit() {
	ChatService = micro.NewService(micro.Name("forum.cli.chat"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))
	ChatService.Init()

	ChatClient = pbc.NewChatServiceClient("forum.service.chat", ChatService.Client())

}
