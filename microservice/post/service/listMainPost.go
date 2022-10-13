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
		Domain:   req.Domain,
		Category: req.Category,
	}

	tag, err := s.Dao.GetTagByContent(req.Tag)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	posts, err := s.Dao.ListMainPost(filter, req.Filter, req.Offset, req.Limit, req.LastId, req.Pagination, req.SearchContent, tag.Id)
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
			CreatorId:     post.CreatorId,
			CreatorName:   post.CreatorName,
			CreatorAvatar: post.CreatorAvatar,
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
