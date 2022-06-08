package service

import (
	"forum-post/dao"
)

// PostService ... 帖子服务
type PostService struct {
	Dao dao.Interface
}

func New(i dao.Interface) *PostService {
	service := new(PostService)
	service.Dao = i
	return service
}
