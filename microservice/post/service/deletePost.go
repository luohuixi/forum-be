package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) DeletePost(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	logger.Info("PostService DeletePost")

	post, err := s.Dao.GetPost(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if post.Re {
		return errno.ServerErr(errno.ErrBadRequest, "this post had been deleted")
	}

	post.Re = true
	if err := post.Save(); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
