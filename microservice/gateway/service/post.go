package service

import (
	pbp "forum-post/proto"
	"forum/pkg/handler"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var PostClient pbp.PostService

func PostInit() {
	service := micro.NewService(micro.Name("forum.cli.post"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))

	service.Init()

	PostClient = pbp.NewPostService("forum.service.post", service.Client())
}
