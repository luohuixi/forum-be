package dao

import pb "forum-post/proto"

type CommentModel struct {
	Id         uint32 `json:"id"`
	TypeId     uint8  `json:"type_id"` // 1 2
	Content    string `json:"content"`
	FatherId   uint32 `json:"father_id"`
	CreateTime string `json:"create_time"`
	Re         bool   `json:"re"`
	CreatorId  uint32 `json:"creator_id"`
	PostId     uint32 `json:"post_id"`
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

func (d *Dao) CreateComment(comment *CommentModel) error {
	return comment.Create()
}

func (d *Dao) ListCommentByPostId(postId uint32) ([]*pb.CommentInfo, error) {
	var comments []*pb.CommentInfo
	err := d.DB.Table("comments").Select("comments.id id, type_id, content, father_id, create_time, creator_id, u.name creator_name, u.avatar creator_avatar, like_num").Joins("join users u on u.id = comments.creator_id").Where("post_id = ? AND re = 0", postId).Find(&comments).Error
	return comments, err
}

func (d *Dao) UpdateCommentInfo(comment *CommentModel) error {
	return comment.Save()
}

func (d *Dao) GetComment(id uint32) (*CommentModel, error) {
	var comment CommentModel
	err := d.DB.Where("id = ? AND re = 0", id).First(&comment).Error
	return &comment, err
}

func (d *Dao) GetCommentNumByPostId(postId uint32) uint32 {
	var count uint32
	d.DB.Model(&CommentModel{}).Where("post_id = ? AND re = 0", postId).Count(&count)
	return count
}
