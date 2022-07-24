package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"forum/util"
)

func (s *PostService) CreateComment(_ context.Context, req *pb.CreateCommentRequest, _ *pb.Response) error {
	logger.Info("PostService CreateComment")

	// check if the FatherId is valid TODO

	data := &dao.CommentModel{
		TypeId:     uint8(req.TypeId),
		Content:    req.Content,
		FatherId:   req.FatherId,
		CreateTime: util.GetCurrentTime(),
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
