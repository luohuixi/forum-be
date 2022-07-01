package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"time"
)

func (s *PostService) CreateComment(ctx context.Context, req *pb.CreateCommentRequest, resp *pb.Response) error {
	logger.Info("PostService CreateComment")

	data := &dao.CommentModel{
		Type:       uint8(req.TypeId),
		Content:    req.Content,
		FatherId:   req.FatherId,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		Re:         false,
		CreatorId:  req.CreatorId,
		PostId:     req.PostId,
	}

	err := s.Dao.CreateComment(data)

	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
