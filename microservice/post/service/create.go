package service

import (
	"context"
	"forum-post/dao"
	errno "forum-post/errno"
	pb "forum-post/proto"
	e "forum/pkg/err"
	"time"
)

func (s *PostService) Create(ctx context.Context, req *pb.CreateRequest, resp *pb.Response) error {
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

	err := s.Dao.Create(data)

	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
