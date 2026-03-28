package main

import (
	"log"

	pb "forum-ocr/proto"
	"forum-ocr/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/identity"
	"forum/pkg/tracer"

	"github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	if err := config.Init("", "FORUM_OCR"); err != nil {
		panic(err)
	}

	t, io, err := tracer.NewTracer(viper.GetString("local_name"), viper.GetString("tracing.jager"))
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	defer logger.SyncLogger()

	opentracing.SetGlobalTracer(t)

	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)

	srv := micro.NewService(
		micro.Name(identity.Prefix()+viper.GetString("local_name")),
		micro.WrapHandler(opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),
		micro.Registry(r),
	)

	srv.Init()

	ocrService, err := service.New()
	if err != nil {
		panic(err)
	}

	if err := pb.RegisterOCRServiceHandler(srv.Server(), ocrService); err != nil {
		panic(err)
	}

	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
