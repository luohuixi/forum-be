package dao

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

type Post2Tag struct {
	Id     uint32
	PostId uint32
	TagId  uint32
}
