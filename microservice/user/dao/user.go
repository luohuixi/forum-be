package dao

import (
	"forum/pkg/constvar"
	"github.com/jinzhu/gorm"
)

type RegisterInfo struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	StudentId string `json:"student_id"`
	Password  string `json:"password"`
	Role      uint32 `json:"role"`
}

// GetUser get a single user by id
func (d *Dao) GetUser(id uint32) (*UserModel, error) {
	user := &UserModel{}
	res := d.DB.Where("id = ?", id).First(user)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, res.Error
}

// GetUserByIds get user by id array
func (d *Dao) GetUserByIds(ids []uint32) ([]*UserModel, error) {
	list := make([]*UserModel, 0)
	if err := d.DB.Where("id IN (?)", ids).Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}

// GetUserByEmail get a user by email.
func (d *Dao) GetUserByEmail(email string) (*UserModel, error) {
	u := &UserModel{}
	err := d.DB.Where("email = ?", email).First(u).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return u, err
}

// GetUserByName get a user by name.
// func GetUserByName(name string) (*UserModel, error) {
// 	u := &UserModel{}
// 	err := m.DB.Self.Where("name = ?", name).First(u).Error
// 	if gorm.IsRecordNotFoundError(err) {
// 		return nil, nil
// 	}
// 	return u, err
// }

// ListUser list users
func (d *Dao) ListUser(offset, limit, lastId uint32, filter *UserModel) ([]*UserModel, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	list := make([]*UserModel, 0)

	query := d.DB.Model(&UserModel{}).Where(filter).Offset(offset).Limit(limit)

	if lastId != 0 {
		query = query.Where("id < ?", lastId).Order("id desc")
	}

	if err := query.Scan(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// GetUserByStudentId get a user by studentId.
func (d *Dao) GetUserByStudentId(studentId string) (*UserModel, error) {
	u := &UserModel{}
	err := d.DB.Where("student_id = ?", studentId).First(u).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return u, err
}

func (d *Dao) RegisterUser(info *RegisterInfo) error {
	// 本地 user 数据库创建用户
	user := &UserModel{
		Name:         info.Name,
		Email:        info.Email,
		StudentID:    info.StudentId,
		HashPassword: generatePasswordHash(info.Password),
		Role:         info.Role,
	}
	return user.Create()
}
