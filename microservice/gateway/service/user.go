package service

import (
	pbu "forum-user/proto"
	"forum/pkg/handler"
	"github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var UserClient pbu.UserService

func UserInit() {

	service := micro.NewService(
		micro.Name("forum.cli.user"),
		micro.Registry(etcd.NewRegistry()),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))

	service.Init()

	UserClient = pbu.NewUserService("forum.service.user", service.Client())
}
