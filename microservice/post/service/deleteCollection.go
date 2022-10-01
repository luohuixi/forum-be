package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"go.uber.org/zap"
)

func (s *PostService) DeleteCollection(_ context.Context, req *pb.Request, _ *pb.Response) error {
	logger.Info("PostService DeleteCollection")

	collection := &dao.CollectionModel{
		Id: req.Id,
	}

	if err := s.Dao.DeleteCollection(collection); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if err := model.DeletePermission(req.UserId, constvar.Collection, req.Id, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	go func() {
		if err := s.Dao.ChangePostScore(req.Id, -constvar.CollectionScore); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), zap.String("cause", err.Error()))
		}
	}()

	return nil
}
