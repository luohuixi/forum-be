package service

import (
	"context"
	pb "forum-chat/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"time"
)

// GetList 获取列表
func (s *ChatService) GetList(ctx context.Context, req *pb.GetListRequest, resp *pb.GetListResponse) error {
	logger.Info("CharService GetList")

	var expiration time.Duration

	ddl, ok := ctx.Deadline()
	if !ok {
		expiration = time.Hour
	} else {
		expiration = ddl.Sub(time.Now())
	}

	// get message of the user from redis
	list, err := s.Dao.GetList(req.UserId, expiration)
	if err != nil {
		return errno.ServerErr(errno.ErrGetRedisList, err.Error())
	}

	select {
	case <-ctx.Done(): // client cancel this request
		if err := s.Dao.Rewrite(req.UserId, list); err != nil {
			return errno.ServerErr(errno.ErrRewriteRedisList, err.Error())
		}
	default:
		if err := s.Dao.CreateHistory(req.UserId, list); err != nil {
			return errno.ServerErr(errno.ErrCreateHistory, err.Error())
		}
		resp.List = list
	}

	return nil
}
