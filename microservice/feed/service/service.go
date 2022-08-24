package service

import (
	"context"
	"forum-feed/dao"
	"forum/pkg/handler"
	opentracingWrapper "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"

	upb "forum-user/proto"

	"github.com/micro/go-micro"
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

var UserService micro.Service
var UserClient upb.UserServiceClient

func UserInit() {
	UserService = micro.NewService(micro.Name("forum.cli.user"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	UserService.Init()

	UserClient = upb.NewUserServiceClient("forum.service.user", UserService.Client())
}

// getInfoFromUserService get user's name and avatar from user-service
func getInfoFromUserService(id uint32) (string, string, error) {
	rsp, err := UserClient.GetProfile(context.TODO(), &upb.GetRequest{Id: id})
	if err != nil {
		return "", "", err
	}

	return rsp.Name, rsp.Avatar, nil
}
