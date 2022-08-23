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
		TypeName:   req.TypeName,
		CategoryId: req.CategoryId,
	}

	posts, err := s.Dao.ListPost(filter, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.List = make([]*pb.Post, len(posts))
	for i, post := range posts {

		// limit max len of post content
		content := post.Content
		if len(content) > 200 {
			content = post.Content[:200]
		}

		isLiked, isFavorite, likeNum, tags, commentNum := s.getPostInfo(post.Id, req.UserId)

		if likeNum != 0 {
			post.LikeNum = likeNum
		}

		resp.List[i] = &pb.Post{
			Id:            post.Id,
			Title:         post.Title,
			Time:          post.LastEditTime,
			CategoryId:    post.CategoryId,
			CreatorId:     post.CreatorId,
			CreatorName:   post.CreatorName,
			CreatorAvatar: post.CreatorAvatar,
			LikeNum:       post.LikeNum,
			Content:       content,
			CommentNum:    commentNum,
			IsLiked:       isLiked,
			IsFavorite:    isFavorite,
			Tags:          tags,
		}
	}

	return nil
}
