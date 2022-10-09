package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) ListLikeByUserId(_ context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListPostPartInfoResponse) error {
	logger.Info("PostService ListLikeByUserId")

	likes, err := s.Dao.ListUserLike(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	var postIds []uint32
	// 暂时忽略对comment的点赞
	for _, like := range likes {
		if like.TypeName == constvar.Post {
			postIds = append(postIds, like.Id)
		}
	}

	posts, err := s.Dao.ListPostInfoByPostIds(postIds, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrListPostInfoByPostIds, err.Error())
	}

	resp.Posts = posts

	return nil
}
