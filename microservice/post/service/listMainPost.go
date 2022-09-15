package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListMainPost(_ context.Context, req *pb.ListMainPostRequest, resp *pb.ListPostResponse) error {
	logger.Info("PostService ListMainPost")

	filter := &dao.PostModel{
		TypeName: req.TypeName,
		Category: req.Category,
	}

	posts, err := s.Dao.ListPost(filter, req.Offset, req.Limit, req.LastId, req.Pagination, req.SearchContent)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Posts = make([]*pb.Post, len(posts))
	for i, post := range posts {

		var content string
		if post.ContentType == "md" {
			content = post.CompiledContent
		} else if post.ContentType == "rtf" {
			content = post.Content
		}

		// limit max len of post content
		if len(content) > 200 {
			content = content[:200]
		}

		isLiked, isCollection, likeNum, tags, commentNum := s.getPostInfo(post.Id, req.UserId)

		if likeNum != 0 {
			post.LikeNum = likeNum
		}

		resp.Posts[i] = &pb.Post{
			Id:            post.Id,
			Title:         post.Title,
			Time:          post.LastEditTime,
			Category:      post.Category,
			CreatorId:     post.CreatorId,
			CreatorName:   post.CreatorName,
			CreatorAvatar: post.CreatorAvatar,
			LikeNum:       post.LikeNum,
			Content:       content,
			CommentNum:    commentNum,
			IsLiked:       isLiked,
			IsCollection:  isCollection,
			Tags:          tags,
			ContentType:   post.ContentType,
		}
	}

	return nil
}
