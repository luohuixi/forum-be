package dao

import (
	"forum/pkg/constvar"
	"gorm.io/gorm"
)

type PostModel struct {
	Id              uint32
	TypeName        string
	Content         string
	Title           string
	CreateTime      string
	Category        string
	Re              bool
	CreatorId       uint32
	LastEditTime    string
	LikeNum         uint32
	ContentType     string
	CompiledContent string
	Summary         string
}

func (PostModel) TableName() string {
	return "posts"
}

// Create ...
func (p *PostModel) Create() error {
	return dao.DB.Create(p).Error
}

// Save ...
func (p *PostModel) Save() error {
	return dao.DB.Save(p).Error
}

func (p *PostModel) Delete() error {
	p.Re = true
	return p.Save()
}

func (p *PostModel) Get(id uint32) error {
	err := dao.DB.Model(p).Where("id = ? AND re = 0", id).First(p).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

type PostInfo struct {
	Id              uint32
	Content         string
	Title           string
	Category        string
	CreatorId       uint32
	LastEditTime    string
	CreatorName     string
	CreatorAvatar   string
	CommentNum      uint32
	LikeNum         uint32
	ContentType     string
	CompiledContent string
	Summary         string
}

func (Dao) CreatePost(post *PostModel) (uint32, error) {
	err := post.Create()
	return post.Id, err
}

func (d *Dao) ListPost(filter *PostModel, offset, limit, lastId uint32, pagination bool, searchContent string) ([]*PostInfo, error) {
	var posts []*PostInfo
	query := d.DB.Table("posts").Select("posts.id id, title, category, compiled_content, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar, content_type, summary").Joins("join users u on u.id = posts.creator_id").Where(filter).Where("posts.re = 0").Order("posts.id desc")

	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		query = query.Offset(int(offset)).Limit(int(limit))

		if lastId != 0 {
			query = query.Where("posts.id < ?", lastId)
		}
	}

	if searchContent != "" {
		// query = query.Where("MATCH (content, title) AGAINST (?)", searchContent) // MySQL 5.7.6 才支持中文全文索引
		query = query.Where("posts.content LIKE ? OR posts.title LIKE ? OR posts.summary LIKE ?", "%"+searchContent+"%", "%"+searchContent+"%", "%"+searchContent+"%")
	}

	err := query.Scan(&posts).Error

	return posts, err
}

func (d *Dao) GetPostInfo(postId uint32) (*PostInfo, error) {
	var post PostInfo
	err := d.DB.Table("posts").Select("posts.id id, title, category, compiled_content, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar, like_num, content_type, summary").Joins("join users u on u.id = posts.creator_id").Where("posts.id = ? AND posts.re = 0", postId).First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &post, err
}

func (Dao) GetPost(id uint32) (*PostModel, error) {
	var item PostModel
	err := item.Get(id)
	return &item, err
}
