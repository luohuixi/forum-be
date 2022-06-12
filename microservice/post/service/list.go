package service

import (
	"context"
	errno "forum-post/errno"
	pb "forum-post/proto"
	e "forum/pkg/err"
)

func (s *PostService) List(ctx context.Context, req *pb.ListRequest, resp *pb.ListResponse) error {
	posts, err := s.Dao.List(uint8(req.TypeId))

	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.List = posts
	return nil
}
