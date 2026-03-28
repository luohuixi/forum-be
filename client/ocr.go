package client

import (
	pbo "forum-ocr/proto"
	"forum/pkg/handler"
	"forum/pkg/identity"

	"github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

var OCRClient pbo.OCRService

func OCRInit() {
	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)
	service := micro.NewService(
		micro.Name("forum.cli.ocr"),
		micro.Registry(r),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	OCRClient = pbo.NewOCRService(identity.Prefix()+"forum.service.ocr", service.Client())
}
