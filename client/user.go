package client

import (
	pbu "forum-user/proto"
	"forum/pkg/identity"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	micro "go-micro.dev/v4"
)

var UserClient pbu.UserService

func UserInit(service micro.Service) {
	UserClient = pbu.NewUserService(identity.Prefix()+"forum.service.user", service.Client())
}
