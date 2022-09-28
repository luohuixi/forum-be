package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"forum/util"
	"gorm.io/gorm"
	"strconv"
)

func (s *PostService) UpdatePostInfo(_ context.Context, req *pb.UpdatePostInfoRequest, _ *pb.Response) error {
	logger.Info("PostService UpdatePostInfo")

	if req.Title == "" || req.Content == "" {
		return errno.ServerErr(errno.ErrBadRequest, "title and content can't be null")
	}

	post, err := s.Dao.GetPost(req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.NotFoundErr(errno.ErrItemNotFound, "post-"+strconv.Itoa(int(req.Id)))
		}

		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if req.Category != post.Category {
		if err := s.Dao.ChangePostCategory(post.TypeName, req.Category, post.Category, req.Id); err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
	}

	post.Title = req.Title
	post.Content = req.Content
	post.CompiledContent = req.CompiledContent
	post.LastEditTime = util.GetCurrentTime()
	post.Category = req.Category
	post.Summary = req.Summary

	if err := post.Save(); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
