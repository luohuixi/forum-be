package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"forum/util"
)

func (s *PostService) CreatePost(_ context.Context, req *pb.CreatePostRequest, _ *pb.Response) error {
	logger.Info("PostService CreatePost")

	if req.MainPostId != 0 {
		post, err := s.Dao.GetPost(req.MainPostId)
		if err != nil {
			return errno.ServerErr(errno.ErrItemNotFound, "main_post cant find")
		}

		if post.MainPostId != 0 {
			return errno.ServerErr(errno.ErrBadRequest, "the main_post_id is not a main_post id")
		}
	}

	data := &dao.PostModel{
		TypeId:       uint8(req.TypeId),
		Content:      req.Content,
		Title:        req.Title,
		CreateTime:   util.GetCurrentTime(),
		Category:     req.Category,
		MainPostId:   req.MainPostId,
		Re:           false,
		CreatorId:    req.UserId,
		LastEditTime: util.GetCurrentTime(),
	}

	if err := s.Dao.CreatePost(data); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
