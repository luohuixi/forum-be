package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"strconv"
)

func (s *PostService) DeleteItem(ctx context.Context, req *pb.Item, resp *pb.Response) error {
	logger.Info("PostService DeleteItem")

	item, err := s.Dao.GetItem(dao.Item{
		Id:     req.Id,
		TypeId: uint8(req.TypeId),
	})
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if item == nil {
		return errno.NotFoundErr(errno.ErrItemNotFound, strconv.Itoa(int(req.Id)))
	}

	if err := item.Delete(); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
