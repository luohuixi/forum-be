package service

import "forum-user/dao"

// UserService ... 用户服务
type UserService struct {
	Dao dao.Interface
}

func New(i dao.Interface) *UserService {
	service := new(UserService)
	service.Dao = i
	return service
}
