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
	SubNum    uint32 // 子评论数
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
	SubNum        uint32
	SubComments   []*CommentInfo
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

// ===================== 游标分页批量获取一级评论 =====================

const commentInfoSelect = `
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
`

func (d *Dao) scanCommentInfos(targetID uint32, targetType string, orderClause string, limit uint32) ([]*CommentInfo, error) {
	var comments []*CommentInfo
	err := d.DB.Table("comments").
		Select(commentInfoSelect).
		Joins("join users u on u.id = comments.creator_id").
		Where("target_id = ? AND target_type = ? AND type_name = ? AND deleted_at = 0 AND is_report = 0",
			targetID, targetType, constvar.FirstLevelComment).
		Order(orderClause).
		Limit(int(limit)).
		Find(&comments).Error
	return comments, err
}

func (d *Dao) scanCommentInfosWithTimeCursor(targetID uint32, targetType string, field string, orderDir string, lastID uint32, lastTime time.Time, limit uint32) ([]*CommentInfo, error) {
	var comments []*CommentInfo
	err := d.DB.Table("comments").
		Select(commentInfoSelect).
		Joins("join users u on u.id = comments.creator_id").
		Where("target_id = ? AND target_type = ? AND type_name = ? AND deleted_at = 0 AND is_report = 0",
			targetID, targetType, constvar.FirstLevelComment).
		Where("("+field+", id) < (?, ?)", lastTime, lastID).
		Order(field + " " + orderDir + ", id " + orderDir).
		Limit(int(limit)).
		Find(&comments).Error
	return comments, err
}

func (d *Dao) scanCommentInfosWithUintCursor(targetID uint32, targetType string, field string, orderDir string, lastID uint32, lastValue uint32, limit uint32) ([]*CommentInfo, error) {
	var comments []*CommentInfo
	err := d.DB.Table("comments").
		Select(commentInfoSelect).
		Joins("join users u on u.id = comments.creator_id").
		Where("target_id = ? AND target_type = ? AND type_name = ? AND deleted_at = 0 AND is_report = 0",
			targetID, targetType, constvar.FirstLevelComment).
		Where("("+field+", id) < (?, ?)", lastValue, lastID).
		Order(field + " " + orderDir + ", id " + orderDir).
		Limit(int(limit)).
		Find(&comments).Error
	return comments, err
}

// ListPrimaryCommentsNewest 按最新排序获取主评论
func (d *Dao) ListPrimaryCommentsNewest(targetID uint32, targetType string, limit uint32) ([]*CommentInfo, error) {
	return d.scanCommentInfos(targetID, targetType, "created_at DESC, id DESC", limit)
}

// ListPrimaryCommentsNewestWithCursor 按最新排序游标分页获取主评论
func (d *Dao) ListPrimaryCommentsNewestWithCursor(targetID uint32, targetType string, lastID uint32, lastTime time.Time, limit uint32) ([]*CommentInfo, error) {
	return d.scanCommentInfosWithTimeCursor(targetID, targetType, "created_at", "DESC", lastID, lastTime, limit)
}

// ListPrimaryCommentsHottest 按最热排序获取主评论
func (d *Dao) ListPrimaryCommentsHottest(targetID uint32, targetType string, limit uint32) ([]*CommentInfo, error) {
	return d.scanCommentInfos(targetID, targetType, "like_num DESC, id DESC", limit)
}

// ListPrimaryCommentsHottestWithCursor 按最热排序游标分页获取主评论
func (d *Dao) ListPrimaryCommentsHottestWithCursor(targetID uint32, targetType string, lastID uint32, lastLikeNum uint32, limit uint32) ([]*CommentInfo, error) {
	return d.scanCommentInfosWithUintCursor(targetID, targetType, "like_num", "DESC", lastID, lastLikeNum, limit)
}

// BatchListCommentsByFatherIDs 批量获取每个 father_id 的前 limit 条子评论
func (d *Dao) BatchListCommentsByFatherIDs(fatherIDs []uint32, limit int) (map[uint32][]*CommentInfo, error) {
	if len(fatherIDs) == 0 {
		return map[uint32][]*CommentInfo{}, nil
	}

	var all []*CommentInfo
	err := d.DB.Table("comments").
		Select(commentInfoSelect).
		Joins("join users u on u.id = comments.creator_id").
		Where("father_id IN ? AND deleted_at = 0 AND is_report = 0", fatherIDs).
		Order("father_id ASC, created_at DESC").
		Find(&all).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint32][]*CommentInfo)
	for _, c := range all {
		list := result[c.FatherId]
		if len(list) >= limit {
			continue
		}
		result[c.FatherId] = append(list, c)
	}
	return result, nil
}

// BatchGetCommentNumByFatherIDs 批量获取多个 father_id 的子评论数
func (d *Dao) BatchGetCommentNumByFatherIDs(fatherIDs []uint32) (map[uint32]uint32, error) {
	if len(fatherIDs) == 0 {
		return map[uint32]uint32{}, nil
	}

	type row struct {
		FatherID uint32
		Count    uint32
	}

	var rows []row
	err := d.DB.Model(&CommentModel{}).
		Select("father_id, COUNT(*) as count").
		Where("father_id IN ? AND deleted_at = 0", fatherIDs).
		Group("father_id").
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint32]uint32)
	for _, r := range rows {
		result[r.FatherID] = r.Count
	}
	return result, nil
}

// ===================== 按 father_id 游标分页获取子评论 =====================

func (d *Dao) scanSubComments(orderClause string, fatherID, limit uint32) ([]*CommentInfo, error) {
	var comments []*CommentInfo
	err := d.DB.Table("comments").
		Select(commentInfoSelect).
		Joins("join users u on u.id = comments.creator_id").
		Where("father_id = ? AND deleted_at = 0 AND is_report = 0", fatherID).
		Order(orderClause).
		Limit(int(limit)).
		Find(&comments).Error
	return comments, err
}

func (d *Dao) scanSubCommentsWithTimeCursor(field, orderDir string, fatherID, lastID uint32, lastTime time.Time, limit uint32) ([]*CommentInfo, error) {
	var comments []*CommentInfo
	err := d.DB.Table("comments").
		Select(commentInfoSelect).
		Joins("join users u on u.id = comments.creator_id").
		Where("father_id = ? AND deleted_at = 0 AND is_report = 0", fatherID).
		Where("("+field+", id) < (?, ?)", lastTime, lastID).
		Order(field + " " + orderDir + ", id " + orderDir).
		Limit(int(limit)).
		Find(&comments).Error
	return comments, err
}

func (d *Dao) scanSubCommentsWithUintCursor(field, orderDir string, fatherID, lastID uint32, lastValue, limit uint32) ([]*CommentInfo, error) {
	var comments []*CommentInfo
	err := d.DB.Table("comments").
		Select(commentInfoSelect).
		Joins("join users u on u.id = comments.creator_id").
		Where("father_id = ? AND deleted_at = 0 AND is_report = 0", fatherID).
		Where("("+field+", id) < (?, ?)", lastValue, lastID).
		Order(field + " " + orderDir + ", id " + orderDir).
		Limit(int(limit)).
		Find(&comments).Error
	return comments, err
}

// ListSubCommentsNewest 按最新排序获取子评论
func (d *Dao) ListSubCommentsNewest(fatherID, limit uint32) ([]*CommentInfo, error) {
	return d.scanSubComments("created_at DESC, id DESC", fatherID, limit)
}

// ListSubCommentsNewestWithCursor 按最新排序游标分页获取子评论
func (d *Dao) ListSubCommentsNewestWithCursor(fatherID, lastID uint32, lastTime time.Time, limit uint32) ([]*CommentInfo, error) {
	return d.scanSubCommentsWithTimeCursor("created_at", "DESC", fatherID, lastID, lastTime, limit)
}

// ListSubCommentsHottest 按最热排序获取子评论
func (d *Dao) ListSubCommentsHottest(fatherID, limit uint32) ([]*CommentInfo, error) {
	return d.scanSubComments("like_num DESC, id DESC", fatherID, limit)
}

// ListSubCommentsHottestWithCursor 按最热排序游标分页获取子评论
func (d *Dao) ListSubCommentsHottestWithCursor(fatherID, lastID uint32, lastLikeNum, limit uint32) ([]*CommentInfo, error) {
	return d.scanSubCommentsWithUintCursor("like_num", "DESC", fatherID, lastID, lastLikeNum, limit)
}
