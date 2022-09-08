package main

import (
	"context"
	pb "forum-feed/proto"
	"forum/pkg/handler"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

func main() {
	service := micro.NewService(micro.Name("forum.cli.feed"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	client := pb.NewFeedService("forum.service.feed", service.Client())

	_, err := client.Delete(context.TODO(), &pb.Request{
		Id: 2,
	})

	panic(err)

	//
	// fmt.Println("post:", post.List[0].Category)
	//
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }
}
