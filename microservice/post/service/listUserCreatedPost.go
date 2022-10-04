package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListUserCreatedPost(_ context.Context, req *pb.Request, resp *pb.ListPostResponse) error {
	logger.Info("PostService ListUserCreatedPost")

	posts, err := s.Dao.ListMyPost(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Posts = make([]*pb.Post, len(posts))
	for i, post := range posts {

		isLiked, isCollection, likeNum, tags, commentNum, collectionNum := s.getPostInfo(post.Id, req.UserId)

		if likeNum != 0 {
			post.LikeNum = likeNum
		}

		resp.Posts[i] = &pb.Post{
			Id:            post.Id,
			Title:         post.Title,
			Time:          post.LastEditTime,
			Category:      post.Category,
			LikeNum:       post.LikeNum,
			CommentNum:    commentNum,
			IsLiked:       isLiked,
			IsCollection:  isCollection,
			Tags:          tags,
			ContentType:   post.ContentType,
			Summary:       post.Summary,
			CollectionNum: collectionNum,
		}
	}

	return nil
}
