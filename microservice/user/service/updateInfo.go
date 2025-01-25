package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

// UpdateInfo ... 更新用户信息
func (s *UserService) UpdateInfo(_ context.Context, req *pb.UpdateInfoRequest, _ *pb.Response) error {
	logger.Info("UserService UpdateInfo")

	user, err := s.Dao.GetUser(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if user == nil {
		return errno.ServerErr(errno.ErrUserNotExisted, "")
	}

	user.Name = req.Info.Name
	user.Avatar = req.Info.AvatarUrl
	user.Signature = req.Info.Signature

	if user.IsPublicCollectionAndLike != req.Info.IsPublicCollectionAndLike {
		if req.Info.IsPublicCollectionAndLike {
			if err := model.AddResourceRole(constvar.CollectionAndLike, user.Id, constvar.CollectionAndLike); err != nil {
				return errno.ServerErr(errno.ErrCasbin, err.Error())
			}

		} else {
			if err := model.DeleteResourceRole(constvar.CollectionAndLike, user.Id, constvar.CollectionAndLike); err != nil {
				return errno.ServerErr(errno.ErrCasbin, err.Error())
			}
		}
		user.IsPublicCollectionAndLike = req.Info.IsPublicCollectionAndLike
	}

	if user.IsPublicFeed != req.Info.IsPublicFeed {
		if req.Info.IsPublicFeed {
			if err := model.AddResourceRole(constvar.Feed, user.Id, constvar.Feed); err != nil {
				return errno.ServerErr(errno.ErrCasbin, err.Error())
			}

		} else {
			if err := model.DeleteResourceRole(constvar.Feed, user.Id, constvar.Feed); err != nil {
				return errno.ServerErr(errno.ErrCasbin, err.Error())
			}
		}

		user.IsPublicFeed = req.Info.IsPublicFeed
	}

	if err := user.Update(); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
