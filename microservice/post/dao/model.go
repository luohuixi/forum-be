package dao

type PostModel struct {
	Id           uint32
	Type         uint8
	Content      string
	Title        string
	CreateTime   string
	Category     string
	Re           bool
	CreatorId    uint32
	LastEditTime string
}

func (u *PostModel) TableName() string {
	return "posts"
}

// Create ... create user
func (u *PostModel) Create() error {
	return dao.DB.Create(u).Error
}

// Save ... save user.
func (u *PostModel) Save() error {
	return dao.DB.Save(u).Error
}
