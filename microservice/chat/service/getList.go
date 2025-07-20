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
	case <-ctx.Done():
		//如果客户端消费失败了,说明要写回去
		if err := s.Dao.Rewrite(req.UserId, list); err != nil {
			return errno.ServerErr(errno.ErrRewriteRedisList, err.Error())
		}
	default:
		// 意思是只要读了就会写到历史记录里面 ?
		if err := s.Dao.CreateHistory(req.UserId, list); err != nil {
			return errno.ServerErr(errno.ErrCreateHistory, err.Error())
		}
		resp.List = list
	}

	return nil
}
