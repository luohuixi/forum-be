package main

import (
	"forum-user/dao"
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	s "forum-user/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	tracer "forum/pkg/tracer"
	"github.com/micro/go-micro"
	"github.com/opentracing/opentracing-go"
	"log"

	// _ "github.com/micro/go-plugins/registry/kubernetes"

	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/spf13/viper"
)

func main() {
	// init config
	if err := config.Init("", "FORUM_USER"); err != nil {
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

	// init oauth-manager and some variables
	auth.InitVar()
	auth.OauthManager.Init()

	srv := micro.NewService(
		micro.Name(viper.GetString("local_name")),
		micro.WrapHandler(
			opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),
	)

	// Init will parse the command line flags.
	srv.Init()

	dao.Init()

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), s.New(dao.GetDao()))

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
