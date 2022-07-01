package main

import (
	"context"
	"forum/config"
	logger "forum/log"
	pb "forum/microservice/post/proto"
	"forum/pkg/handler"
	tracer "forum/pkg/tracer"
	"github.com/micro/go-micro"
	"github.com/opentracing/opentracing-go"
	"log"

	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/spf13/viper"
)

func main() {
	// init config
	if err := config.Init("", "FORUM_post"); err != nil {
		panic(err)
	}

	t, io, err := tracer.NewTracer(viper.GetString("local_name"), viper.GetString("tracing.jager"))
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	defer logger.SyncLogger()

	// set var t to Global Tracer (opentracing single instance mode)
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(micro.Name("forum.cli.post"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	client := pb.NewPostServiceClient("forum.service.post", service.Client())

	client.CreatePost(context.Background(), &pb.CreatePostRequest{
		UserId:   0,
		Content:  "",
		TypeId:   0,
		Title:    "",
		Category: "",
	})
}
