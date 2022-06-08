package dao

type PostModel struct {
	Id       uint32
	Type     uint8
	Context  string
	Title    string
	Time     string
	Category string
	Re       bool
	Creator  uint
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
