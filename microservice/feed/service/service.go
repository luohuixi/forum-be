package service

import (
	"context"
	"forum-feed/dao"
	"forum/pkg/handler"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"

	upb "forum-user/proto"

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

var UserClient upb.UserService

func UserInit() {
	service := micro.NewService(micro.Name("forum.cli.user"),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),
	)

	service.Init()

	UserClient = upb.NewUserService("forum.service.user", service.Client())
}

// getInfoFromUserService get user's name and avatar from user-service
func getInfoFromUserService(id uint32) (string, string, error) {
	rsp, err := UserClient.GetProfile(context.TODO(), &upb.GetRequest{Id: id})
	if err != nil {
		return "", "", err
	}

	return rsp.Name, rsp.Avatar, nil
}
