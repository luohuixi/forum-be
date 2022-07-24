package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"
)

func (s *PostService) ListPost(_ context.Context, req *pb.ListPostRequest, resp *pb.ListPostResponse) error {
	logger.Info("PostService ListPost")

	typeId, err := strconv.Atoi(req.TypeId)
	if err != nil {
		return errno.ServerErr(errno.ErrBadRequest, err.Error())
	}

	filter := &dao.PostModel{TypeId: uint8(typeId)}

	if req.Category != "" {
		filter.Category = req.Category
	}
	// TODO: limit offset
	posts, err := s.Dao.ListPost(filter)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.List = make([]*pb.Post, len(posts))
	for i, post := range posts {
		commentNum := s.Dao.GetCommentNumByPostId(post.Id)

		item := dao.Item{
			Id:     post.Id,
			TypeId: constvar.Post,
		}

		isLiked, err := s.Dao.IsUserHadLike(req.UserId, item)
		if err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}

		likeNum, err := s.Dao.GetLikedNum(item)
		if err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}

		// limit max len of post content
		content := post.Content
		if len(content) > 200 {
			content = post.Content[:200]
		}

		resp.List[i] = &pb.Post{
			Id:            post.Id,
			Title:         post.Title,
			Time:          post.LastEditTime,
			Category:      post.Category,
			CreatorId:     post.CreatorId,
			CreatorName:   post.CreatorName,
			CreatorAvatar: post.CreatorAvatar,
			Content:       content,
			CommentNum:    commentNum,
			LikeNum:       uint32(likeNum),
			IsLiked:       isLiked,
		}
	}

	return nil
}
