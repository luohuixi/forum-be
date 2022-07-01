package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) GetPost(ctx context.Context, req *pb.Request, resp *pb.Post) error {
	logger.Info("PostService GetPost")

	// ok, err := s.Dao.Enforce(userId, typeId, constvar.Read)
	post, err := s.Dao.GetPostInfo(req.Id)

	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	comments, err := s.Dao.ListCommentByPostId(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	likeNum, err := s.Dao.GetLikedNum(dao.Item{
		Id:     req.Id,
		TypeId: 1,
	})
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	resp = &pb.Post{
		Id:            post.Id,
		Content:       post.Content,
		Title:         post.Title,
		Time:          post.LastEditTime,
		Category:      post.Category,
		CreatorId:     post.CreatorId,
		CreatorAvatar: post.CreatorAvatar,
		CreatorName:   post.CreatorName,
		Comments:      comments,
		CommentNum:    uint32(len(comments)),
		LikeNum:       uint32(likeNum),
	}
	return nil
}
