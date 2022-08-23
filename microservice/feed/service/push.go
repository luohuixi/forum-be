package service

import (
	"context"
	"encoding/json"
	"forum/pkg/errno"
	"forum/util"

	"forum-feed/dao"
	pb "forum-feed/proto"
	logger "forum/log"
)

// Push ... 异步新增feed
func (s *FeedService) Push(_ context.Context, req *pb.PushRequest, _ *pb.Response) error {
	logger.Info("FeedService Push")

	// get username and avatar by userId from user-service
	username, avatar, err := getInfoFromUserService(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRPC, err.Error())
	}

	feed := &dao.FeedModel{
		UserId:           req.UserId,
		UserName:         username,
		UserAvatar:       avatar,
		Action:           req.Action,
		SourceTypeName:   req.Source.TypeName,
		SourceObjectName: req.Source.Name,
		SourceObjectId:   req.Source.Id,
		CreateTime:       util.GetCurrentTime(),
	}

	msg, _ := json.Marshal(feed)

	if err := dao.PublishMsg(msg); err != nil {
		return errno.ServerErr(errno.ErrPublishMsg, err.Error())
	}

	return nil
}
