package service

import (
	pbc "forum-chat/proto"
	handler "forum/pkg/handler"

	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
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
