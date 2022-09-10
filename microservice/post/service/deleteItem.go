package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) DeleteItem(_ context.Context, req *pb.DeleteItemRequest, _ *pb.Response) error {
	logger.Info("PostService DeleteItem")

	err := s.Dao.DeleteItem(dao.Item{
		Id:       req.Id,
		TypeName: req.TypeName,
	})
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if err := model.DeletePermission(req.UserId, req.TypeName, req.Id, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	return nil
}
