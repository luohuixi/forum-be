package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/util"
	"go.uber.org/zap"
	"strconv"
)

func (s *PostService) GetComment(_ context.Context, req *pb.Request, resp *pb.CommentInfo) error {
	logger.Info("PostService GetComment")

	comment, err := s.Dao.GetCommentInfo(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if comment == nil {
		return errno.NotFoundErr(errno.ErrItemNotFound, "comment-"+strconv.Itoa(int(req.Id)))
	}

	likeNum, err := s.Dao.GetLikedNum(dao.Item{
		Id:       req.Id,
		TypeName: constvar.Comment,
	})
	if err != nil {
		logger.Error(err.Error(), zap.Error(errno.ErrRedis))
	}

	resp.LikeNum = comment.LikeNum
	if likeNum != 0 {
		resp.LikeNum = uint32(likeNum)
	}
	resp.TypeName = comment.TypeName
	resp.Id = comment.Id
	resp.Content = comment.Content
	resp.CreateTime = util.FormatString(comment.CreateTime)
	resp.CreatorId = comment.CreatorId
	resp.CreatorAvatar = comment.CreatorAvatar
	resp.CreatorName = comment.CreatorName

	return nil
}
