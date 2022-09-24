package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListCollection(_ context.Context, req *pb.UserIdRequest, resp *pb.ListCollectionsResponse) error {
	logger.Info("PostService ListCollections")

	collections, err := s.Dao.ListCollectionByUserId(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	for _, collection := range collections {

		commentNum, err := s.Dao.GetCommentNumByPostId(collection.PostId)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}

		collection.CommentNum = commentNum
	}

	resp.Collections = collections

	return nil
}
