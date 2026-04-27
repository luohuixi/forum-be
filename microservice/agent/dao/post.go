package dao

type PostModel struct {
	Id      uint32
	Content string
	Title   string
}

// todo：change
func (d *Dao) GetPostByTime(time string) (*[]PostModel, error) {
	var allPosts []PostModel
	if err := d.DB.Table("posts").Where("create_time >= ?", time).Find(&allPosts).Error; err != nil {
		return nil, err
	}

	return &allPosts, nil
}
