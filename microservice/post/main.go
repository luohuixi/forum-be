package main

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	"forum-post/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/tracer"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/server"
	"log"

	_ "github.com/micro/go-plugins/registry/kubernetes"

	"github.com/spf13/viper"
)

// logWrapper is a handler wrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[wrapper] server request: %v", req.Endpoint())
		err := fn(ctx, req, rsp)
		return err
	}
}

func main() {
	// init config
	if err := config.Init("", "FORUM_POST"); err != nil {
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

	srv := micro.NewService(
		micro.Name(viper.GetString("local_name")),
		// micro.WrapHandler(
		// 	opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		// ),
		micro.WrapHandler(logWrapper),
	)

	// Init will parse the command line flags.
	srv.Init()

	dao.Init()

	// Register handler
	pb.RegisterPostServiceHandler(srv.Server(), service.New(dao.GetDao()))

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
