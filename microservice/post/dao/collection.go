package dao

import (
	pb "forum-post/proto"
)

type CollectionModel struct {
	Id         uint32
	UserId     uint32
	CreateTime string
	PostId     uint32
}

func (CollectionModel) TableName() string {
	return "collections"
}

// Create ...
func (c *CollectionModel) Create() error {
	return dao.DB.Create(c).Error
}

func (c *CollectionModel) Delete() error {
	return dao.DB.Delete(c).Error
}

func (Dao) CreateCollection(collection *CollectionModel) (uint32, error) {
	err := collection.Create()
	return collection.Id, err
}

func (Dao) DeleteCollection(collection *CollectionModel) error {
	return collection.Delete()
}

func (d *Dao) ListCollectionByUserId(userId uint32) ([]*pb.Collection, error) {
	var collections []*pb.Collection
	err := d.DB.Table("collections").Select("collections.id id, post_id, title, content, collections.create_time time, p.creator_id, u.name creator_name, u.avatar creator_avatar").Joins("join posts p on p.id = collections.post_id").Joins("join users u on u.id = p.creator_id").Where("user_id = ?", userId).Find(&collections).Error
	return collections, err
}

func (d *Dao) IsUserCollectionPost(userId uint32, postId uint32) (bool, error) {
	var count int64
	err := d.DB.Table("collections").Where("user_id = ? AND post_id = ?", userId, postId).Count(&count).Error
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (d *Dao) GetCollectionNumByPostId(postId uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&CollectionModel{}).Where("post_id = ? AND re = 0", postId).Count(&count).Error
	return uint32(count), err
}
