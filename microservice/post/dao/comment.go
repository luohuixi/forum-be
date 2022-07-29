package dao

import (
	pb "forum-post/proto"
	"github.com/jinzhu/gorm"
)

type CommentModel struct {
	Id         uint32 `json:"id"`
	TypeId     uint8  `json:"type_id"` // 1 : 一级; 2 : 二级
	Content    string `json:"content"`
	FatherId   uint32 `json:"father_id"`
	CreateTime string `json:"create_time"`
	Re         bool   `json:"re"`
	CreatorId  uint32 `json:"creator_id"`
	PostId     uint32 `json:"post_id"`
	LikeNum    uint32 `json:"like_num"`
}

func (c *CommentModel) TableName() string {
	return "comments"
}

// Create ...
func (c *CommentModel) Create() error {
	return dao.DB.Create(c).Error
}

// Save ...
func (c *CommentModel) Save() error {
	return dao.DB.Save(c).Error
}

func (c *CommentModel) Get(id uint32) error {
	err := dao.DB.Model(c).Where("id = ? AND re = 0", id).First(c).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (c *CommentModel) Delete() error {
	c.Re = true
	return c.Save()
}

type CommentInfo struct {
	Id            uint32 `json:"id"`
	TypeId        uint8  `json:"type_id"` // 1 2
	Content       string `json:"content"`
	FatherId      uint32 `json:"father_id"`
	CreateTime    string `json:"create_time"`
	CreatorId     uint32 `json:"creator_id"`
	PostId        uint32 `json:"post_id"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	LikeNum       uint32 `json:"like_num"`
	// TODO
}

func (d *Dao) CreateComment(comment *CommentModel) error {
	return comment.Create()
}

func (d *Dao) ListCommentByPostId(postId uint32) ([]*pb.CommentInfo, error) {
	var comments []*pb.CommentInfo
	err := d.DB.Table("comments").Select("comments.id id, type_id, content, father_id, create_time, creator_id, u.name creator_name, u.avatar creator_avatar, like_num").Joins("join users u on u.id = comments.creator_id").Where("post_id = ? AND comments.re = 0", postId).Find(&comments).Error
	return comments, err
}

func (d *Dao) GetComment(id uint32) (*CommentModel, error) {
	var comment CommentModel
	err := d.DB.Where("id = ? AND re = 0", id).First(&comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (d *Dao) GetCommentInfo(commentId uint32) (*CommentInfo, error) {
	var comment CommentInfo // TODO
	err := d.DB.Table("comments").Select("posts.id id, title, category, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar, like_num").Joins("join users u on u.id = posts.creator_id").Where("posts.id = ? AND comments.re = 0", commentId).First(&comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (d *Dao) GetCommentNumByPostId(postId uint32) uint32 {
	var count uint32
	d.DB.Model(&CommentModel{}).Where("post_id = ? AND re = 0", postId).Count(&count)
	return count
}
