package main

import (
	"context"
	"fmt"
	pb "forum-post/proto"
	"forum/pkg/handler"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

func main() {
	service := micro.NewService(micro.Name("forum.cli.post"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	client := pb.NewPostService("forum.service.post", service.Client())

	resp, err := client.ListReport(context.TODO(), &pb.ListReportRequest{})

	fmt.Println("----- resp: ", resp, " -----")

	panic(err)

	//
	// fmt.Println("post:", post.List[0].Category)
	//
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }
}
