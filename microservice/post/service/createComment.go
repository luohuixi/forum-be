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
	"go.uber.org/zap"
)

func (s *PostService) CreateComment(_ context.Context, req *pb.CreateCommentRequest, resp *pb.CreateResponse) error {
	logger.Info("PostService CreateComment")

	// check if the FatherId is valid
	switch req.TypeName {
	case constvar.SubPost:
		post, err := s.Dao.GetPostInfo(req.FatherId)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if post == nil {
			return errno.ServerErr(errno.ErrBadRequest, "the post not found")
		}

		resp.TargetUserId = post.Id
		resp.TargetUserId = post.CreatorId

	case constvar.FirstLevelComment, constvar.SecondLevelComment:
		comment, err := s.Dao.GetComment(req.FatherId)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if comment == nil {
			return errno.ServerErr(errno.ErrBadRequest, "the comment not found")
		}

		if (req.TypeName == constvar.FirstLevelComment && comment.TypeName != constvar.SubPost) || (req.TypeName == constvar.SecondLevelComment && comment.TypeName != constvar.FirstLevelComment) {
			return errno.ServerErr(errno.ErrBadRequest, "type_name of father not legal")
		}

		resp.TargetUserId = comment.Id
		resp.TargetUserId = comment.CreatorId

	default:
		return errno.ServerErr(errno.ErrBadRequest, "type_name not legal")
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

	go func() {
		if err := s.Dao.ChangePostScore(req.PostId, constvar.CommentScore); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), zap.String("cause", err.Error()))
		}
	}()

	return nil
}
