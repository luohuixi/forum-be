package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	// "github.com/micro/micro/v3/service/logger"
	"strconv"
)

func (s *PostService) ListPost(ctx context.Context, req *pb.ListPostRequest, resp *pb.ListPostResponse) error {
	logger.Info("PostService ListPost")

	id, err := strconv.Atoi(req.TypeId)
	if err != nil {
		return errno.ServerErr(errno.ErrBadRequest, err.Error())
	}

	var posts []*dao.PostInfo
	if req.Category == "" {
		posts, err = s.Dao.ListPost(uint8(id))
	} else {
		posts, err = s.Dao.ListPostByCategory(uint8(id), req.Category)
	}
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	for _, post := range posts {
		commentNum := s.Dao.GetCommentNumByPostId(post.Id)

		likeNum, err := s.Dao.GetLikedNum(dao.Item{
			Id:     post.Id,
			TypeId: 1,
		})
		if err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}

		// limit max len of post content
		content := post.Content
		if len(content) > 200 {
			content = post.Content[:200]
		}

		postInfo := &pb.Post{
			Id:            post.Id,
			Title:         post.Title,
			Time:          post.LastEditTime,
			Category:      post.Category,
			CreatorId:     post.CreatorId,
			CreatorName:   post.CreatorName,
			Content:       content,
			CreatorAvatar: post.CreatorAvatar,
			CommentNum:    commentNum,
			LikeNum:       uint32(likeNum),
		}
		resp.List = append(resp.List, postInfo)
	}

	return nil
}
