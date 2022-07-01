package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) DeletePost(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	logger.Info("PostService DeletePost")

	comment, err := s.Dao.GetComment(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	comment.Re = true
	if err := s.Dao.UpdateCommentInfo(comment); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
