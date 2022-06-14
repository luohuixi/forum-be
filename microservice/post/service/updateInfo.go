package service

import (
	"context"
	errno "forum-post/errno"
	pb "forum-post/proto"
	logger "forum/log"
	e "forum/pkg/err"
	"time"
)

func (s *PostService) UpdateInfo(ctx context.Context, req *pb.UpdateInfoRequest, resp *pb.Response) error {
	logger.Info("PostService UpdateInfo")

	post, err := s.Dao.Get(req.Id)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	post.Title = req.Title
	post.Content = req.Content
	post.LastEditTime = time.Now().Format("2006-01-02 15:04:05")
	post.Category = req.Category

	if err := s.Dao.UpdateInfo(post); err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
