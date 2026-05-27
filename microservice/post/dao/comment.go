package dao

import (
	"forum/pkg/constvar"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// todo 索引

type CommentModel struct {
	Id         uint32
	TargetID   uint32
	TargetType string
	CreatedAt  time.Time
	DeletedAt  soft_delete.DeletedAt `gorm:"softDelete:nano"`

	TypeName  string // constvar.FirstLevelComment or constvar.SecondLevelComment
	Content   string
	FatherId  uint32
	CreatorId uint32
	LikeNum   uint32
	ImgUrl    string
	IsReport  bool
}

func (CommentModel) TableName() string {
	return "comments"
}

// Create ...
func (c *CommentModel) Create() error {
	return dao.DB.Create(c).Error
}

func (c *CommentModel) Delete(tx *gorm.DB) error {
	return tx.Delete(c).Error
}

func (c *CommentModel) Get(id uint32) error {
	return dao.DB.Where("id = ? AND deleted_at = 0", id).First(c).Error
}

// Save ...
func (c *CommentModel) Save(tx ...*gorm.DB) error {
	db := dao.DB
	if len(tx) == 1 {
		db = tx[0]
	}

	return db.Save(c).Error
}

func (c *CommentModel) BeReported() error {
	c.IsReport = true
	return c.Save()
}

func (c *CommentModel) CancelReported() error {
	c.IsReport = false
	return c.Save()
}

type CommentInfo struct {
	Id            uint32
	TargetID      uint32
	TargetType    string
	CreatedAt     time.Time
	TypeName      string
	Content       string
	FatherId      uint32
	CreatorId     uint32
	LikeNum       uint32
	ImgUrl        string
	CreatorName   string
	CreatorAvatar string
}

func (d *Dao) CreateComment(comment *CommentModel, tx ...*gorm.DB) (uint32, error) {
	db := d.getDB(tx...)
	if comment.TargetType == "" {
		comment.TargetType = constvar.Post
	}
	err := db.Create(comment).Error
	return comment.Id, err
}

func (d *Dao) ListCommentByPostId(postId uint32) ([]*CommentInfo, error) {
	var comments []*CommentInfo

	err := d.DB.Table("comments").
		Select(`
			comments.id id,
			type_name,
			content,
			father_id,
			created_at,
			creator_id,
			u.name creator_name,
			u.avatar creator_avatar,
			like_num,
			img_url
		`).
		Joins("join users u on u.id = comments.creator_id").
		Where(
			"target_id = ? AND (target_type = ? OR target_type = '' OR target_type IS NULL) AND deleted_at = 0 AND is_report = 0",
			postId,
			constvar.Post,
		).
		Find(&comments).Error

	return comments, err
}

func (d *Dao) ListCommentByTarget(targetId uint32, targetType string) ([]*CommentInfo, error) {
	var comments []*CommentInfo

	err := d.DB.Table("comments").
		Select(`
			comments.id id,
			type_name,
			content,
			father_id,
			created_at,
			creator_id,
			u.name creator_name,
			u.avatar creator_avatar,
			like_num,
			img_url
		`).
		Joins("join users u on u.id = comments.creator_id").
		Where(
			"target_id = ? AND target_type = ? AND deleted_at = 0 AND is_report = 0",
			targetId,
			targetType,
		).
		Find(&comments).Error

	return comments, err
}

func (d *Dao) GetComment(id uint32, tx ...*gorm.DB) (*CommentModel, error) {
	db := d.getDB(tx...)
	var comment CommentModel

	err := db.
		Where("id = ? AND deleted_at = 0", id).
		First(&comment).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (d *Dao) GetCommentInfo(commentId uint32) (*CommentInfo, error) {
	var comment CommentInfo

	err := d.DB.Table("comments").
		Select(`
			comments.id id,
			type_name,
			content,
			father_id,
			created_at,
			comments.creator_id,
			target_id,
			target_type,
			u.name creator_name,
			u.avatar creator_avatar,
			like_num,
			img_url
		`).
		Joins("join users u on u.id = comments.creator_id").
		Where("comments.id = ? AND deleted_at = 0", commentId).
		First(&comment).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (d *Dao) GetCommentNumByPostId(postId uint32) (uint32, error) {
	var count int64

	err := d.DB.Model(&CommentModel{}).
		Where(
			"target_id = ? AND (target_type = ? OR target_type = '' OR target_type IS NULL) AND deleted_at = 0",
			postId,
			constvar.Post,
		).
		Count(&count).Error

	return uint32(count), err
}

func (d *Dao) GetCommentNumByTarget(targetID uint32, targetType string) (uint32, error) {
	var count int64

	err := d.DB.Model(&CommentModel{}).
		Where(
			"target_id = ? AND target_type = ? AND deleted_at = 0",
			targetID,
			targetType,
		).
		Count(&count).Error

	return uint32(count), err
}

// BatchListCommentsByTargets 批量获取每个 target 的前 limit 条评论
func (d *Dao) BatchListCommentsByTargets(targetIDs []uint32, targetType string, limit int) (map[uint32][]*CommentInfo, error) {
	if len(targetIDs) == 0 {
		return map[uint32][]*CommentInfo{}, nil
	}

	var all []*CommentInfo
	err := d.DB.Table("comments").
		Select(`
			comments.id id,
			type_name,
			content,
			father_id,
			created_at,
			creator_id,
			target_id,
			u.name creator_name,
			u.avatar creator_avatar,
			like_num,
			img_url
		`).
		Joins("join users u on u.id = comments.creator_id").
		Where(
			"target_id IN ? AND target_type = ? AND deleted_at = 0 AND is_report = 0",
			targetIDs,
			targetType,
		).
		Order("target_id ASC, created_at DESC").
		Find(&all).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint32][]*CommentInfo)
	for _, c := range all {
		list := result[c.TargetID]
		if len(list) >= limit {
			continue
		}
		result[c.TargetID] = append(list, c)
	}
	return result, nil
}

// BatchGetCommentNumByTargets 批量获取多个 target 的评论数
func (d *Dao) BatchGetCommentNumByTargets(targetIDs []uint32, targetType string) (map[uint32]uint32, error) {
	if len(targetIDs) == 0 {
		return map[uint32]uint32{}, nil
	}

	type row struct {
		TargetID uint32
		Count    uint32
	}

	var rows []row
	err := d.DB.Model(&CommentModel{}).
		Select("target_id, COUNT(*) as count").
		Where("target_id IN ? AND target_type = ? AND deleted_at = 0", targetIDs, targetType).
		Group("target_id").
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint32]uint32)
	for _, r := range rows {
		result[r.TargetID] = r.Count
	}
	return result, nil
}
