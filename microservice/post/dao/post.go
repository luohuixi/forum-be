package dao

import (
	"forum/pkg/constvar"
	"github.com/jinzhu/gorm"
)

type PostModel struct {
	Id           uint32 `json:"id"`
	TypeId       uint8  `json:"type_id"`
	Content      string `json:"content"`
	Title        string `json:"title"`
	CreateTime   string `json:"create_time"`
	CategoryId   uint32 `json:"category_id"`
	Re           bool   `json:"re"`
	CreatorId    uint32 `json:"creator_id"`
	LastEditTime string `json:"last_edit_time"`
	LikeNum      uint32 `json:"like_num"`
	MainPostId   uint32 `json:"main_post_id"`
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
	Id            uint32 `json:"id"`
	Content       string `json:"content"`
	Title         string `json:"title"`
	CategoryId    uint32 `json:"category_id"`
	CreatorId     uint32 `json:"creator_id"`
	LastEditTime  string `json:"last_edit_time"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	CommentNum    uint32 `json:"comment_num"`
	LikeNum       uint32 `json:"like_num"`
}

func (Dao) CreatePost(post *PostModel) (uint32, error) {
	err := post.Create()
	return post.Id, err
}

func (d *Dao) ListPost(filter *PostModel, offset, limit, lastID uint32, pagination bool) ([]*PostInfo, error) {
	var posts []*PostInfo
	query := d.DB.Table("posts").Select("posts.id id, title, category_id, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar").Joins("join users u on u.id = posts.creator_id").Where(filter).Where("posts.re = 0")

	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		query = query.Offset(offset).Limit(limit)

		if lastID != 0 {
			query = query.Where("projects.id < ?", lastID)
		}
	}

	err := query.Scan(&posts).Error

	return posts, err
}

func (d *Dao) GetPostInfo(postId uint32) (*PostInfo, error) {
	var post PostInfo
	err := d.DB.Table("posts").Select("posts.id id, title, category_id, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar, like_num").Joins("join users u on u.id = posts.creator_id").Where("posts.id = ? AND posts.re = 0", postId).First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &post, err
}

func (Dao) GetPost(id uint32) (*PostModel, error) {
	item := &PostModel{}
	err := item.Get(id)
	return item, err
}

func (d *Dao) IsUserFavoritePost(userId uint32, id uint32) (bool, error) {
	// TODO
	return false, nil
}
