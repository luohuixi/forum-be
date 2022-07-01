package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	e "forum/pkg/err"
	errno "forum/pkg/err"
)

func (s *PostService) CreateLike(ctx context.Context, req *pb.LikeRequest, resp *pb.Response) error {
	logger.Info("PostService CreateLike")

	item := dao.Item{
		Id:     req.TargetId,
		TypeId: uint8(req.TypeId),
	}

	ok, err := s.Dao.IsUserHadLike(req.UserId, item)
	if err != nil {
		return e.ServerErr(errno.ErrRedis, err.Error())
	}
	if ok {
		return e.ServerErr(errno.ErrBadRequest, "已点赞")
	}

	err = s.Dao.AddLike(req.UserId, item)
	if err != nil {
		return e.ServerErr(errno.ErrRedis, err.Error())
	}

	return nil
}
