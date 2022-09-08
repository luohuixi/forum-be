package service

import (
	pbp "forum-post/proto"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	// opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	// "github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var PostClient pbp.PostService

func PostInit() {
	service := micro.NewService(micro.Name("forum.cli.post")) // micro.WrapClient(
	// 	opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
	// ),
	// micro.WrapCall(handler.ClientErrorHandlerWrapper()))
	// wrap the client
	// micro.WrapClient(logWrap),

	service.Init()

	PostClient = pbp.NewPostService("forum.service.post", service.Client())
}
