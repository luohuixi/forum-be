package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) CreateOrRemoveLike(_ context.Context, req *pb.LikeRequest, _ *pb.Response) error {
	logger.Info("PostService CreateOrRemoveLike")

	var score int

	item := dao.Item{
		Id:       req.Item.TargetId,
		TypeName: req.Item.TypeName,
	}

	ok, err := s.Dao.IsUserHadLike(req.UserId, item)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	if ok {
		err = s.Dao.RemoveLike(req.UserId, item)
		score = -constvar.LikeScore
	} else {
		err = s.Dao.AddLike(req.UserId, item)
		score = constvar.LikeScore
	}
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	if req.Item.TypeName == constvar.Post {
		go func() {
			if err := s.Dao.AddChangeRecord(req.Item.TargetId); err != nil {
				logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
			}

			if err := s.Dao.ChangePostScore(req.Item.TargetId, score); err != nil {
				logger.Error(errno.ErrChangeScore.Error(), logger.String(err.Error()))
			}
		}()
	}

	return nil
}
