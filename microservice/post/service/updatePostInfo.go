package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"time"
)

func (s *PostService) UpdatePostInfo(ctx context.Context, req *pb.UpdateInfoRequest, resp *pb.Response) error {
	logger.Info("PostService UpdatePostInfo")

	post, err := s.Dao.GetPost(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	post.Title = req.Title
	post.Content = req.Content
	post.LastEditTime = time.Now().Format("2006-01-02 15:04:05")
	post.Category = req.Category

	if err := s.Dao.UpdatePostInfo(post); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
