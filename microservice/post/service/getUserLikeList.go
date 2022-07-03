package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListLikeByUserId(ctx context.Context, req *pb.Request, resp *pb.ListLikeResponse) error {
	logger.Info("PostService ListLikeByUserId")

	likes, err := s.Dao.ListUserLike(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	for _, like := range likes {
		resp.List = append(resp.List, &pb.LikeItem{
			TargetId: like.Id,
			TypeId:   uint32(like.TypeId),
		})
	}

	return nil
}
