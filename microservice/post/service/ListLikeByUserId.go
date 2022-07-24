package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListLikeByUserId(_ context.Context, req *pb.UserIdRequest, resp *pb.ListLikeResponse) error {
	logger.Info("PostService ListLikeByUserId")

	likes, err := s.Dao.ListUserLike(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	resp.List = make([]*pb.LikeItem, len(likes))
	for i, like := range likes {
		resp.List[i] = &pb.LikeItem{
			TargetId: like.Id,
			TypeId:   uint32(like.TypeId),
		}
	}

	return nil
}
