package service

import (
	pbf "forum-feed/proto"
	"forum/pkg/handler"

	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
)

var FeedService micro.Service
var FeedClient pbf.FeedServiceClient

func FeedInit() {
	FeedService = micro.NewService(micro.Name("forum.cli.feed"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))
	FeedService.Init()

	FeedClient = pbf.NewFeedServiceClient("forum.service.feed", FeedService.Client())

}
