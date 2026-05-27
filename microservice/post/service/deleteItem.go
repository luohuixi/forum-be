package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"gorm.io/gorm"
)

func (s *PostService) DeleteItem(_ context.Context, req *pb.DeleteItemRequest, _ *pb.Response) error {
	logger.Info("PostService DeleteItem")

	var err error

	switch req.TypeName {
	case constvar.Post:
		err = s.Dao.DeletePost(req.Id)
	case constvar.Comment:
		err = s.deleteComment(req.Id)
	case constvar.QualityPost:
		err = s.Dao.ChangeQualityPost(req.Id, false)
	case constvar.SipScore:
		err = s.deleteSipScore(req.Id)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "wrong TypeName")
	}

	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if req.TypeName != constvar.QualityPost {
		if err := model.DeletePermission(req.UserId, req.TypeName, req.Id, constvar.Write); err != nil {
			return errno.ServerErr(errno.ErrCasbin, err.Error())
		}
	}

	return nil
}

func (s *PostService) deleteSipScore(id uint32) error {
	// 删除 主体
	if err := s.Dao.DeleteSipScore(id); err != nil {
		return err
	}

	// todo 消息队列中比较好
	// 删除 sipScore tags
	// 删除 sipScoreEntry, sipScoreEntryComment, sipScoreEntryReview
	// 删除 sipScore collection
	return nil
}

func (s *PostService) deleteComment(id uint32) error {
	// 1. 获取对应评论
	comment, err := s.Dao.GetComment(id)
	if err != nil {
		return err
	}
	if comment == nil {
		return nil // 评论已不存在，幂等处理
	}

	// 2. 查看对应的targetType
	switch comment.TargetType {
	// 3. post 则根据之前的直接写（删除评论 + 减少 post score）
	case constvar.Post, "":
		if err := s.Dao.DeleteComment(id); err != nil {
			return err
		}
		return s.Dao.ChangePostScore(comment.TargetID, -constvar.CommentScore)

	// 4. SipScoreEntryCommentRating 则需要减小对应 SipScoreEntryCommentRating 的 commentNum
	case constvar.SipScoreEntryCommentRating:
		rating, err := s.Dao.GetSipScoreEntryCommentRatingByID(comment.TargetID)
		if err != nil {
			return err
		}
		if rating == nil {
			return s.Dao.DeleteComment(id)
		}

		return s.Dao.Transaction(func(tx *gorm.DB) error {
			if err := s.Dao.DeleteComment(id, tx); err != nil {
				return err
			}
			return s.Dao.DecrSipScoreEntryCommentRatingCommentNum(rating.SipScoreID, rating.EntryID, rating.ID, tx)
		})

	default:
		return s.Dao.DeleteComment(id)
	}
}
