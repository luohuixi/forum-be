package service

import (
	pbu "forum-user/proto"
	"forum/pkg/handler"

	_ "github.com/micro/go-plugins/registry/kubernetes"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var UserService micro.Service
var UserClient pbu.UserServiceClient

func UserInit() {
	UserService = micro.NewService(micro.Name("forum.cli.user"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))
	UserService.Init()

	UserClient = pbu.NewUserServiceClient("forum.service.user", UserService.Client())

}
