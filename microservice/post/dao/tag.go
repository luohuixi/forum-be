package dao

import (
	logger "forum/log"
	"forum/pkg/errno"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type TagModel struct {
	Id      uint32
	Content string
}

func (TagModel) TableName() string {
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
	content, err := d.getTagContentById(strconv.Itoa(int(id)))
	if err != nil {
		return tag, err
	}
	if content != "" {
		tag.Content = content
		return tag, nil
	}

	// 从redis缓存中未命中则在数据库找
	if err := dao.DB.Model(tag).First(tag).Error; err != nil {
		return nil, err
	}

	if err := dao.addTag(tag.Id, tag.Content); err != nil {
		logger.Error(err.Error(), zap.Error(errno.ErrRedis))
	}

	return tag, nil
}

func (d *Dao) GetTagByContent(content string) (*TagModel, error) {
	tag := &TagModel{
		Content: content,
	}

	if content == "" {
		return tag, nil
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

	if err := dao.addTag(tag.Id, tag.Content); err != nil {
		logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
	}

	return tag, err
}

func (d *Dao) getTagContentById(id string) (string, error) {
	content, err := d.Redis.Get("tag:id:" + id).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	err = d.Redis.Expire("tag:id:"+id, 10*24*time.Hour).Err()
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

func (d *Dao) addTag(id uint32, content string) error {
	pipe := d.Redis.TxPipeline()

	pipe.Set("tag:id:"+strconv.Itoa(int(id)), content, 10*24*time.Hour)
	pipe.Set("tag:content:"+content, id, 10*24*time.Hour)

	_, err := pipe.Exec()
	return err
}

type Post2TagModel struct {
	Id     uint32
	PostId uint32
	TagId  uint32
}

func (Post2TagModel) TableName() string {
	return "post2tags"
}

// Create ...
func (p *Post2TagModel) Create() error {
	return dao.DB.Create(p).Error
}

func (Dao) CreatePost2Tag(item Post2TagModel) error {
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

func (d *Dao) AddTagToSortedSet(tagId uint32) error {
	return d.Redis.ZIncrBy("tags", 1, strconv.Itoa(int(tagId))).Err() // FIXME test
	// d.Redis.ZAdd("tags", redis.Z{Score: 1, Member: strconv.Itoa(int(tagId))}).Err()
}

func (d *Dao) ListPopularTags() ([]string, error) {
	// 降序
	ids, err := d.Redis.ZRevRange("tags", 0, 9).Result()
	if err != nil {
		return nil, err
	}

	tags := make([]string, len(ids))
	for i, id := range ids {
		tags[i], err = d.getTagContentById(id)
		if err != nil {
			return nil, err
		}
	}

	return tags, nil
}
