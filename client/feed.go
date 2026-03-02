package client

import (
	pbf "forum-feed/proto"
	"forum/pkg/identity"

	micro "go-micro.dev/v4"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
)

var FeedClient pbf.FeedService

func FeedInit(service micro.Service) {
	FeedClient = pbf.NewFeedService(identity.Prefix()+"forum.service.feed", service.Client())
}
