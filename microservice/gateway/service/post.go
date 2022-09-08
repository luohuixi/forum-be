package service

import (
	pbp "forum-post/proto"
	"forum/pkg/handler"

	_ "github.com/micro/go-plugins/registry/kubernetes"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var PostService micro.Service
var PostClient pbp.PostServiceClient

func PostInit() {
	PostService = micro.NewService(micro.Name("forum.cli.post"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))
	PostService.Init()

	PostClient = pbp.NewPostServiceClient("forum.service.post", PostService.Client())

}
