package dao

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

func (d Dao) DeleteCollection(collection *CollectionModel) error {
	return d.DB.Table("collections").Where("post_id = ? AND user_id = ?", collection.PostId, collection.UserId).Delete(collection).Error
}

func (d *Dao) ListCollectionByUserId(userId uint32) ([]uint32, error) {
	var postIds []uint32
	err := d.DB.Table("collections").Select("post_id").Where("user_id = ?", userId).Find(&postIds).Error
	return postIds, err
}

func (d *Dao) IsUserCollectionPost(userId uint32, postId uint32) (bool, error) {
	var count int64
	if err := d.DB.Table("collections").Where("user_id = ? AND post_id = ?", userId, postId).Count(&count).Error; err != nil {
		return false, err
	}

	return count != 0, nil
}

func (d *Dao) GetCollectionNumByPostId(postId uint32) (uint32, error) {
	var count int64
	err := d.DB.Model(&CollectionModel{}).Where("post_id = ?", postId).Count(&count).Error
	return uint32(count), err
}
