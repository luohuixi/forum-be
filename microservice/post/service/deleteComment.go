package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	e "forum/pkg/err"
	errno "forum/pkg/err"
)

func (s *PostService) DeleteComment(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	logger.Info("PostService DeleteComment")

	post, err := s.Dao.GetPost(req.Id)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	post.Re = true
	if err := s.Dao.UpdatePostInfo(post); err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
