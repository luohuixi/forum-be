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

func (s *PostService) CreateComment(_ context.Context, req *pb.CreateCommentRequest, _ *pb.Response) error {
	logger.Info("PostService CreateComment")

	// check if the FatherId is valid
	if req.TypeName == constvar.FirstLevelComment {
		post, err := s.Dao.GetPost(req.FatherId)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if post == nil {
			return errno.ServerErr(errno.ErrBadRequest, "the post not found")
		}
	} else if req.TypeName == constvar.SecondLevelComment {
		comment, err := s.Dao.GetComment(req.FatherId)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if comment == nil {
			return errno.ServerErr(errno.ErrBadRequest, "the comment not found")
		}
	} else {
		return errno.ServerErr(errno.ErrBadRequest, "TypeName must be "+constvar.FirstLevelComment+" or "+constvar.SecondLevelComment)
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

	commentId, err := s.Dao.CreateComment(data)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if err := model.AddPolicy(req.CreatorId, constvar.Comment, commentId, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	if err := model.AddRole(constvar.Comment, commentId, constvar.Post+":"+req.TypeName); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	return nil
}
