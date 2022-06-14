package service

import (
	"context"
	errno "forum-post/errno"
	pb "forum-post/proto"
	logger "forum/log"
	e "forum/pkg/err"
	// "github.com/micro/micro/v3/service/logger"
	"strconv"
)

func (s *PostService) List(ctx context.Context, req *pb.ListRequest, resp *pb.ListResponse) error {
	logger.Info("PostService List")

	id, err := strconv.Atoi(req.TypeId)
	if err != nil {
		return e.ServerErr(errno.ErrBadRequest, err.Error())
	}
	posts, err := s.Dao.List(uint8(id))

	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.List = posts
	return nil
}
