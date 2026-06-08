package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) CreateFeedback(_ context.Context, req *pb.CreateFeedbackRequest, _ *pb.Response) error {
	logger.Info("PostService CreateFeedback")

	if req.GetContent() == "" {
		return errno.ServerErr(errno.ErrBadRequest, "content required")
	}

	feedback := &dao.FeedbackModel{
		UserID:   req.GetUserId(),
		Category: req.GetCategory(),
		Content:  req.GetContent(),
		Contact:  req.GetContact(),
		ImgURL:   req.GetImgUrl(),
	}
	if err := s.Dao.CreateFeedback(feedback); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	return nil
}
