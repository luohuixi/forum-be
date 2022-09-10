package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/util"
)

func (s *PostService) CreateCollection(_ context.Context, req *pb.Request, resp *pb.CreateResponse) error {
	logger.Info("PostService CreateCollection")

	collection := &dao.CollectionModel{
		CreateTime: util.GetCurrentTime(),
		UserId:     req.UserId,
		PostId:     req.Id,
	}

	collectionId, err := s.Dao.CreateCollection(collection)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if err := model.AddPolicy(req.UserId, constvar.Collection, collectionId, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	resp.Id = collectionId

	return nil
}
