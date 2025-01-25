package main

import (
	"forum-post/dao"
	pb "forum-post/proto"
	"forum-post/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/tracer"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
	"log"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"

	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/spf13/viper"
)

func init() {
	service.UserInit()
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
		micro.WrapHandler(
			opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),
		micro.Registry(etcd.NewRegistry()),
	)

	// Init will parse the command line flags.
	srv.Init()

	dao.Init()

	// Register handler
	if err := pb.RegisterPostServiceHandler(srv.Server(), service.New(dao.GetDao())); err != nil {
		panic(err)
	}

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
