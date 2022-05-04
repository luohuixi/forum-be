package service

import (
	"context"
	"forum-chat/errno"
	pb "forum-chat/proto"
	e "forum/pkg/err"
	"log"
	"time"
)

// GetList 获取列表
func (s *ChatService) GetList(ctx context.Context, req *pb.GetListRequest, resp *pb.GetListResponse) error {
	log.Println("CharService.GetList", req.UserId)

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
		return e.ServerErr(errno.ErrGetRedisList, err.Error())
	}

	select {
	case <-ctx.Done(): // client cancel this request
		err := s.Dao.Rewrite(req.UserId, list)
		if err != nil {
			return e.ServerErr(errno.ErrRewriteRedisList, err.Error())
		}
	default:
		resp.List = list
	}

	return nil
}
