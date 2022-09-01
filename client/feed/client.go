package main

import (
	"context"
	pb "forum-feed/proto"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/tracer"
	"github.com/micro/go-micro"
	"github.com/opentracing/opentracing-go"
	"log"

	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
)

func main() {
	t, io, err := tracer.NewTracer("forum.service.feed", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	defer logger.SyncLogger()

	opentracing.SetGlobalTracer(t)

	service := micro.NewService(micro.Name("forum.cli.feed"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	client := pb.NewFeedServiceClient("forum.service.feed", service.Client())

	_, err = client.Delete(context.TODO(), &pb.Request{
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
