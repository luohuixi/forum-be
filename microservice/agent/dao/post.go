package dao

type PostModel struct {
	Id      uint32
	Content string
	Title   string
}

func (d *Dao) GetPostById(id uint32) (*PostModel, error) {
	var post PostModel
	if err := d.DB.Table("posts").Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
