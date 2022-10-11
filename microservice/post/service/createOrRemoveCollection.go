package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/util"
)

func (s *PostService) CreateOrRemoveCollection(_ context.Context, req *pb.Request, resp *pb.CreateCollectionResponse) error {
	logger.Info("PostService CreateOrRemoveCollection")

	var score int

	collection := &dao.CollectionModel{
		CreateTime: util.GetCurrentTime(),
		UserId:     req.UserId,
		PostId:     req.Id,
	}

	isCollection, err := s.Dao.IsUserCollectionPost(req.UserId, req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if isCollection {
		err = s.Dao.DeleteCollection(collection)

		score = -constvar.CollectionScore
	} else {
		resp.Id, err = s.Dao.CreateCollection(collection)

		score = constvar.CollectionScore
	}
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	go func() {
		if err := s.Dao.ChangePostScore(req.Id, score); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), logger.String(err.Error()))
		}
	}()

	post, err := s.Dao.GetPost(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.UserId = post.CreatorId
	resp.Content = post.Title
	resp.TypeName = post.TypeName

	return nil
}
