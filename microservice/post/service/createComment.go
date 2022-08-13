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

func (s *PostService) CreateComment(_ context.Context, req *pb.CreateCommentRequest, _ *pb.Response) error {
	logger.Info("PostService CreateComment")

	// check if the FatherId is valid
	if req.TypeName == constvar.FirstLevelComment {
		post, err := s.Dao.GetPost(req.FatherId)
		if err != nil {
			errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if post == nil {
			errno.ServerErr(errno.ErrBadRequest, "the post not found")
		}
	} else if req.TypeName == constvar.SecondLevelComment {
		comment, err := s.Dao.GetComment(req.FatherId)
		if err != nil {
			errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if comment == nil {
			errno.ServerErr(errno.ErrBadRequest, "the comment not found")
		}
	} else {
		return errno.ServerErr(errno.ErrBadRequest, "TypeName should be first-level or second-level")
	}

	data := &dao.CommentModel{
		TypeName:   req.TypeName,
		Content:    req.Content,
		FatherId:   req.FatherId,
		CreateTime: util.GetCurrentTime(),
		Re:         false,
		CreatorId:  req.CreatorId,
		PostId:     req.PostId,
	}

	if err := s.Dao.CreateComment(data); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
