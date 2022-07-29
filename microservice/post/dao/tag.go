package dao

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type TagModel struct {
	Id      uint32
	Content string
}

func (t *TagModel) TableName() string {
	return "tags"
}

// Create ...
func (t *TagModel) Create() error {
	return dao.DB.Create(t).Error
}

func (d *Dao) GetTagById(id uint32) (*TagModel, error) {
	tag := &TagModel{
		Id: id,
	}
	content, err := d.getTagContentById(id)
	if err != nil {
		return tag, err
	}
	if content != "" {
		tag.Content = content
		return tag, nil
	}

	// 从redis缓存中未命中则在数据库找
	err = dao.DB.Model(tag).Where("id = ?", id).First(tag).Error

	return tag, err
}

func (d *Dao) GetTagByContent(content string) (*TagModel, error) {
	tag := &TagModel{
		Content: content,
	}
	id, err := d.getTagIdByContent(content)
	if err != nil {
		return tag, err
	}
	if id != 0 {
		tag.Id = uint32(id)
		return tag, nil
	}

	// 从redis缓存中未命中则在数据库找
	err = dao.DB.Model(tag).Where("content = ?", content).First(tag).Error
	if err == gorm.ErrRecordNotFound {
		// 在数据库未找到则新建
		err = tag.Create()
	}

	return tag, err
}

func (d *Dao) getTagContentById(id uint32) (string, error) {
	content, err := d.Redis.Get("tag:id:" + strconv.Itoa(int(id))).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	err = d.Redis.Expire("tag:id:"+strconv.Itoa(int(id)), 10*24*time.Hour).Err()
	return content, err
}

func (d *Dao) getTagIdByContent(content string) (int, error) {
	id, err := d.Redis.Get("tag:content:" + content).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	err = d.Redis.Expire("tag:content:"+content, 10*24*time.Hour).Err()
	return id, err
}

func (d *Dao) AddTag(id uint32, content string) error {
	pipe := d.Redis.TxPipeline()

	pipe.Set("tag:id:"+strconv.Itoa(int(id)), content, 10*24*time.Hour).Err()
	pipe.Set("tag:content:"+content, id, 10*24*time.Hour)

	_, err := pipe.Exec()
	return err
}

type Post2TagModel struct {
	Id     uint32
	PostId uint32
	TagId  uint32
}

func (p *Post2TagModel) TableName() string {
	return "post2tags"
}

// Create ...
func (p *Post2TagModel) Create() error {
	return dao.DB.Create(p).Error
}

func (d *Dao) CreatePost2Tag(item Post2TagModel) error {
	return item.Create()
}

func (d *Dao) ListTagsByPostId(postId uint32) ([]string, error) {
	var list []*Post2TagModel
	err := d.DB.Table("post2tags").Where("post_id = ?", postId).Find(&list).Error
	if err != nil {
		return nil, err
	}

	contents := make([]string, len(list))
	for i, item := range list {
		tag, err := d.GetTagById(item.TagId)
		if err != nil {
			return nil, err
		}
		contents[i] = tag.Content
	}

	return contents, nil
}
