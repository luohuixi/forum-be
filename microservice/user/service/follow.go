package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *UserService) ToggleFollow(_ context.Context, req *pb.FollowRequest, resp *pb.FollowResponse) error {
	logger.Info("UserService ToggleFollow")

	userID := req.GetUserId()
	targetUserID := req.GetTargetUserId()
	if userID == 0 || targetUserID == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "user_id and target_user_id required")
	}
	if userID == targetUserID {
		return errno.ServerErr(errno.ErrBadRequest, "cannot follow yourself")
	}

	target, err := s.Dao.GetUser(targetUserID)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if target == nil {
		return errno.ServerErr(errno.ErrUserNotExisted, "")
	}

	resp.IsFollowing, err = s.Dao.ToggleFollow(userID, targetUserID)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	resp.FollowingCount, err = s.Dao.CountFollowing(userID)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	resp.FollowerCount, err = s.Dao.CountFollowers(targetUserID)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	return nil
}
