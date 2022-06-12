package dao

import (
	pb "forum-post/proto"
)

type PostInfo struct {
	Type     uint8
	Content  string
	Title    string
	Time     string
	Category string
	Creator  uint32
}

func (d *Dao) Create(post *PostModel) error {
	return post.Create()
}

func (d *Dao) List(typeId uint8) ([]*pb.Post, error) {
	var posts []*pb.Post
	err := d.DB.Table("posts").Find(posts).Error
	return posts, err
}

func (d *Dao) UpdateInfo(post *PostModel) error {
	return post.Save()
}

func (d *Dao) Get(postId uint32) (*PostModel, error) {
	var post *PostModel
	err := d.DB.First(post).Error
	return post, err
}
