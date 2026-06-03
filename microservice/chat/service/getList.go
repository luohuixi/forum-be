package service

import (
	"context"
	"fmt"
	pb "forum-chat/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"time"

	"go.uber.org/zap"
)

// GetList 获取列表
func (s *ChatService) GetList(ctx context.Context, req *pb.GetListRequest, resp *pb.GetListResponse) error {
	logger.Info(fmt.Sprintf("CharService GetList user(%d)", req.UserId))

	var expiration time.Duration
	ddl, ok := ctx.Deadline()
	if !ok {
		expiration = time.Hour
	} else {
		expiration = ddl.Sub(time.Now())
	}

	// get message of the user from redis
	list, err := s.Dao.GetList(req.UserId, expiration, req.Wait)
	if err != nil {
		return errno.ServerErr(errno.ErrGetRedisList, err.Error())
	}

	if ctx.Err() != nil {
		go s.Rewrite(req.UserId, list)
		resp.List = list
		return nil
	}

	if len(list) > 0 {
		if err := s.Dao.CreateHistory(req.UserId, list); err != nil {
			go s.Rewrite(req.UserId, list)
			return errno.ServerErr(errno.ErrCreateHistory, err.Error())
		}
	}

	resp.List = list

	return nil
}

func (s *ChatService) Rewrite(userid uint32, list []string) {
	if err := s.Dao.Rewrite(userid, list); err != nil {
		rewriteErr := errno.ServerErr(errno.ErrRewriteRedisList, err.Error())
		logger.Error("Rewrite Failed",
			zap.String("cause", rewriteErr.Error()),
			zap.Strings("source", list),
		)
	}
}

func (s *ChatService) CreateHistory(userid uint32, list []string) {
	if err := s.Dao.CreateHistory(userid, list); err != nil {
		historyErr := errno.ServerErr(errno.ErrCreateHistory, err.Error())
		logger.Error("Create History Failed",
			zap.String("cause", historyErr.Error()),
			zap.Strings("source", list),
		)
	}
}
