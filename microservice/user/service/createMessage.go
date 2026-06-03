package service

import (
	"context"
	"encoding/json"
	"strconv"

	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"

	"github.com/google/uuid"
)

// CreateMessage ... 创建用户消息
func (s *UserService) CreateMessage(_ context.Context, req *pb.CreateMessageRequest, _ *pb.Response) error {
	logger.Info("UserService CreateMessage")

	if err := s.Dao.CreateMessage(0, req.Message); err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	return nil
}

func (s *UserService) CreatePrivateMessage(_ context.Context, req *pb.CreatePrivateMessageRequest, _ *pb.Response) error {
	logger.Info("UserService CreatePrivateMessage")

	uid := uuid.New().String()
	message, _ := json.Marshal(map[string]string{
		"id":              uid,
		"send_user_id":    strconv.Itoa(int(req.SendUserId)),
		"content":         req.Content,
		"post_id":         strconv.Itoa(int(req.PostId)),
		"comment_id":      strconv.Itoa(int(req.CommentId)),
		"type":            req.Type,
		"post_title":      req.PostTitle,
		"comment_content": req.CommentContent,
	})

	var err error
	if req.Type == "like" || req.Type == "collection" {
		err = s.Dao.CreateOrUpdateInteractionMessage(req.ReceiveUserId, string(message))
	} else {
		err = s.Dao.CreateMessage(req.ReceiveUserId, string(message))
	}
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	return nil
}
