package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"time"
)

func (s *PostService) CreatePost(ctx context.Context, req *pb.CreatePostRequest, resp *pb.Response) error {
	logger.Info("PostService CreatePost")

	data := &dao.PostModel{
		Type:         uint8(req.TypeId),
		Content:      req.Content,
		Title:        req.Title,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		Category:     req.Category,
		Re:           false,
		CreatorId:    req.UserId,
		LastEditTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	err := s.Dao.CreatePost(data)

	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
