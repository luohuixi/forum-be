package dao

type PostModel struct {
	Id           uint32 `json:"id"`
	Type         uint8  `json:"type"`
	Content      string `json:"content"`
	Title        string `json:"title"`
	CreateTime   string `json:"create_time"`
	Category     string `json:"category"`
	Re           bool   `json:"re"`
	CreatorId    uint32 `json:"creator_id"`
	LastEditTime string `json:"last_edit_time"`
}

func (p *PostModel) TableName() string {
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

type PostInfo struct {
	Id            uint32 `json:"id"`
	Content       string `json:"content"`
	Title         string `json:"title"`
	Category      string `json:"category"`
	CreatorId     uint32 `json:"creator_id"`
	LastEditTime  string `json:"last_edit_time"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	CommentNum    uint32 `json:"comment_num"`
}

func (d *Dao) CreatePost(post *PostModel) error {
	return post.Create()
}

func (d *Dao) ListPost(typeId uint8) ([]*PostInfo, error) {
	var posts []*PostInfo
	err := d.DB.Table("posts p").Select("p.id id", "title", "category", "content", "last_edit_time", "creator_id", "u.name creator_name", "u.avatar creator_avatar").Joins("join users u on u.id = p.creator_id").Where("type = ? AND re = 0", typeId).Find(posts).Error
	return posts, err
}

func (d *Dao) UpdatePostInfo(post *PostModel) error {
	return post.Save()
}

func (d *Dao) GetPost(postId uint32) (*PostModel, error) {
	var post *PostModel
	err := d.DB.Where("id = ? AND re = 0", postId).First(post).Error
	return post, err
}

func (d *Dao) GetPostInfo(postId uint32) (*PostInfo, error) {
	var post *PostInfo
	err := d.DB.Table("posts p").Select("p.id id", "title", "category", "content", "last_edit_time", "creator_id", "u.name creator_name", "u.avatar creator_avatar").Joins("join users u on u.id = p.creator_id").Where("id = ? AND re = 0", postId).First(post).Error
	return post, err
}
