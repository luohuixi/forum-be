package service

import (
	"context"
	"forum-user/dao"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *UserService) ListFollowing(_ context.Context, req *pb.FollowListRequest, resp *pb.FollowListResponse) error {
	logger.Info("UserService ListFollowing")
	return s.listFollowUsers(req, resp, true)
}

func (s *UserService) ListFollowers(_ context.Context, req *pb.FollowListRequest, resp *pb.FollowListResponse) error {
	logger.Info("UserService ListFollowers")
	return s.listFollowUsers(req, resp, false)
}

func (s *UserService) listFollowUsers(req *pb.FollowListRequest, resp *pb.FollowListResponse, following bool) error {
	if req.GetUserId() == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "user_id required")
	}

	limit := req.GetLimit()
	if limit == 0 || limit > 100 {
		limit = 50
	}

	var (
		users []*dao.FollowListUser
		err   error
	)
	if following {
		users, err = s.Dao.ListFollowing(req.GetUserId(), limit, req.GetOffset())
	} else {
		users, err = s.Dao.ListFollowers(req.GetUserId(), limit, req.GetOffset())
	}
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if len(users) == 0 {
		resp.Users = []*pb.FollowListUser{}
		return nil
	}

	userIDs := make([]uint32, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.Id)
	}

	followingCounts, err := s.Dao.BatchCountFollowing(userIDs)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	followerCounts, err := s.Dao.BatchCountFollowers(userIDs)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	followingMap, err := s.Dao.BatchIsFollowing(req.GetViewerId(), userIDs)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Users = make([]*pb.FollowListUser, 0, len(users))
	for _, user := range users {
		resp.Users = append(resp.Users, &pb.FollowListUser{
			Id:             user.Id,
			Name:           user.Name,
			Avatar:         user.Avatar,
			Role:           user.Role,
			Signature:      user.Signature,
			FollowingCount: followingCounts[user.Id],
			FollowerCount:  followerCounts[user.Id],
			IsFollowing:    followingMap[user.Id],
		})
	}
	return nil
}
