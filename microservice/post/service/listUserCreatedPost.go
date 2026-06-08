package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) ListUserCreatedPost(ctx context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListPostPartInfoResponse) error {
	logger.Info("PostService ListUserCreatedPost")

	targetUserId := req.GetTargetUserId()
	if targetUserId == 0 {
		targetUserId = req.GetUserId()
	}
	viewerUserId := req.GetUserId()
	if viewerUserId == 0 {
		viewerUserId = targetUserId
	}

	postIds, err := s.Dao.ListUserCreatedPost(targetUserId, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	var filter = &dao.PostModel{}

	domain, err := s.GetUserDomain(ctx, viewerUserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRPC, err.Error())
	}
	if domain == constvar.NormalDomain {
		filter.Domain = domain
	}

	posts, err := s.Dao.ListPostInfoByPostIds(postIds, filter, 0, 0, 0, false)
	if err != nil {
		return errno.ServerErr(errno.ErrListPostInfoByPostIds, err.Error())
	}

	resp.Posts = posts

	return nil
}
