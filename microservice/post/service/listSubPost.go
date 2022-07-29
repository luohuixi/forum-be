package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) ListSubPost(_ context.Context, req *pb.ListSubPostRequest, resp *pb.ListPostResponse) error {
	logger.Info("PostService ListSubPost")

	filter := &dao.PostModel{
		TypeId:     uint8(req.TypeId),
		MainPostId: req.MainPostId,
	}

	posts, err := s.Dao.ListPost(filter, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.List = make([]*pb.Post, len(posts))
	for i, post := range posts {
		// TODO comments etc.

		// limit max len of post content
		content := post.Content
		if len(content) > 200 {
			content = post.Content[:200]
		}

		resp.List[i] = &pb.Post{
			Id:            post.Id,
			Title:         post.Title,
			Time:          post.LastEditTime,
			CategoryId:    post.CategoryId,
			CreatorId:     post.CreatorId,
			CreatorName:   post.CreatorName,
			CreatorAvatar: post.CreatorAvatar,
			Content:       content,
		}

		resp.List[i].CommentNum = s.Dao.GetCommentNumByPostId(post.Id)

		item := dao.Item{
			Id:     post.Id,
			TypeId: constvar.Post,
		}

		resp.List[i].IsLiked, err = s.Dao.IsUserHadLike(req.UserId, item)
		if err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}

		likeNum, err := s.Dao.GetLikedNum(item)
		if err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
		resp.List[i].LikeNum = uint32(likeNum)
	}

	return nil
}
