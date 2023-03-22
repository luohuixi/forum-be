package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) ListLikeByUserId(_ context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListPostPartInfoResponse) error {
	logger.Info("PostService ListLikeByUserId")

	likes, err := s.Dao.ListUserLike(req.TargetUserId)
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

	var filter = &dao.PostModel{}

	domain, err := s.GetUserDomain(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRPC, err.Error())
	}
	if domain == constvar.NormalDomain {
		filter.Domain = domain
	}

	posts, err := s.Dao.ListPostInfoByPostIds(postIds, filter, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrListPostInfoByPostIds, err.Error())
	}

	resp.Posts = posts

	return nil
}
