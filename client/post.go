package client

import (
	pbp "forum-post/proto"
	"forum/pkg/identity"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	micro "go-micro.dev/v4"
)

var PostClient pbp.PostService

func PostInit(service micro.Service) {
	PostClient = pbp.NewPostService(identity.Prefix()+"forum.service.post", service.Client())
}
