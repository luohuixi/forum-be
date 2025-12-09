package service

import (
	"forum-feed/dao"
	"forum/pkg/handler"
	"github.com/go-micro/plugins/v4/registry/etcd"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go-micro.dev/v4/registry"

	pbu "forum-user/proto"

	micro "go-micro.dev/v4"
)

// FeedService ... 动态服务
type FeedService struct {
	Dao dao.Interface
}

func New(i dao.Interface) *FeedService {
	service := new(FeedService)
	service.Dao = i
	return service
}

var UserClient pbu.UserService

func UserInit() {
	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)
	service := micro.NewService(micro.Name("forum.cli.user"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.Registry(r),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	UserClient = pbu.NewUserService("forum.service.user", service.Client())
}
