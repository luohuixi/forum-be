package dao

import (
	"forum/pkg/constvar"
	"gorm.io/gorm"
	"strconv"
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

func (p *PostModel) Update() error {
	err := dao.DB.Table("posts").Where("id = ?", p.Id).Updates(map[string]interface{}{
		"title":            p.Title,
		"content":          p.Content,
		"compiled_content": p.CompiledContent,
		"last_edit_time":   p.LastEditTime,
		"category":         p.Category,
		"summary":          p.Summary,
	}).Error

	return err
}

func (p *PostModel) Delete() error {
	p.Re = true
	return p.Save()
}

func (p *PostModel) Get(id uint32) error {
	return dao.DB.Model(p).Where("id = ? AND re = 0", id).First(p).Error
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

func (d *Dao) ListMainPost(filter *PostModel, offset, limit, lastId uint32, pagination bool, searchContent string, tagId uint32) ([]*PostInfo, error) {
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

	if tagId != 0 {
		var postIds []uint32
		if err := d.DB.Table("post2tags").Select("post_id").Distinct("post_id").Where("tag_id = ?", tagId).Find(&postIds).Error; err != nil {
			return nil, err
		}

		query.Where("posts.id IN ?", postIds)
	}

	if searchContent != "" {
		// query = query.Where("MATCH (content, title) AGAINST (?)", searchContent) // MySQL 5.7.6 才支持中文全文索引
		query = query.Where("posts.content LIKE ? OR posts.title LIKE ? OR posts.summary LIKE ?", "%"+searchContent+"%", "%"+searchContent+"%", "%"+searchContent+"%")
	}

	err := query.Scan(&posts).Error

	return posts, err
}

func (d *Dao) ListMyPost(creatorId uint32) ([]*PostInfo, error) {
	var posts []*PostInfo
	err := d.DB.Table("posts").Select("posts.id id, title, category, compiled_content, content, last_edit_time, creator_id, u.name creator_name, u.avatar creator_avatar, content_type, summary").Joins("join users u on u.id = posts.creator_id").Where("creator_id = ?", creatorId).Where("posts.re = 0").Order("posts.id desc").Scan(&posts).Error

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

func (d Dao) ChangePostScore(postId uint32, score int) error {
	post, err := d.GetPost(postId)
	if err != nil {
		return err
	}

	if err := d.Redis.ZIncrBy("hot:"+post.TypeName+":"+post.Category, float64(score), strconv.Itoa(int(postId))).Err(); err != nil {
		return err
	}
	return d.Redis.ZIncrBy("hot:"+post.TypeName, float64(score), strconv.Itoa(int(postId))).Err()
}

func (d Dao) ChangePostCategory(typeName, newCategory, oldCategory string, postId uint32) error {
	score, err := d.Redis.ZScore("hot"+typeName+":"+oldCategory, strconv.Itoa(int(postId))).Result()
	if err != nil {
		return err
	}

	if err := d.Redis.ZRem("hot"+typeName+":"+oldCategory, strconv.Itoa(int(postId))).Err(); err != nil {
		return err
	}

	return d.Redis.ZIncrBy("hot:"+typeName+":"+newCategory, score, strconv.Itoa(int(postId))).Err()
}

func (d Dao) ListHotPost(typeName, category string, offset, limit uint32, pagination bool) ([]*PostInfo, error) {
	key := "hot:" + typeName
	if category != "" {
		key += ":" + category
	}

	var start int64
	var end int64 = -1
	if pagination {
		if limit == 0 {
			limit = constvar.DefaultLimit
		}

		start = int64(offset)
		end = int64(offset + limit)
	}

	result, err := d.Redis.ZRevRange(key, start, end).Result()
	if err != nil {
		return nil, err
	}

	list := make([]*PostInfo, len(result))
	for i, r := range result {
		data, err := strconv.Atoi(r)
		if err != nil {
			return nil, err
		}

		list[i], err = d.GetPostInfo(uint32(data))
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
