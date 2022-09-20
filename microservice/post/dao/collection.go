package dao

import (
	pb "forum-post/proto"
	"gorm.io/gorm"
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
	err := d.DB.Table("collections").Select("collections.id id, post_id, title, content, create_time time, creator_id, u.name creator_name, u.avatar creator_avatar").Joins("join users u on u.id = collections.creator_id").Where("user_id = ?", userId).Find(&collections).Error
	return collections, err
}

func (d *Dao) IsUserCollectionPost(userId uint32, postId uint32) (bool, error) {
	err := d.DB.Table("collections").Where("user_id = ? AND post_id = ?", userId, postId).First(&CollectionModel{}).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
