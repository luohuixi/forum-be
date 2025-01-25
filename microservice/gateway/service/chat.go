package service

import (
	pbc "forum-chat/proto"
	"forum/pkg/handler"
	"github.com/go-micro/plugins/v4/registry/etcd"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var ChatClient pbc.ChatService

func ChatInit() {
	service := micro.NewService(
		micro.Name("forum.cli.chat"),
		micro.Registry(etcd.NewRegistry()),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))

	service.Init()

	ChatClient = pbc.NewChatService("forum.service.chat", service.Client())
}
