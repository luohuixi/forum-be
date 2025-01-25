package service

import (
	pbf "forum-feed/proto"
	"forum/pkg/handler"
	"github.com/go-micro/plugins/v4/registry/etcd"

	micro "go-micro.dev/v4"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
)

var FeedClient pbf.FeedService

func FeedInit() {
	service := micro.NewService(
		micro.Name("forum.cli.feed"),
		micro.Registry(etcd.NewRegistry()),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))

	service.Init()

	FeedClient = pbf.NewFeedService("forum.service.feed", service.Client())
}
