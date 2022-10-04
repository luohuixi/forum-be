package dao

import (
	"forum/pkg/constvar"
	"gorm.io/gorm"
)

type RegisterInfo struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	StudentId string `json:"student_id"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

// GetUser get a single user by id
func (d *Dao) GetUser(id uint32) (*UserModel, error) {
	user := &UserModel{}
	err := d.DB.Where("id = ? AND re = 0", id).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

// GetUserByIds get user by id array
func (d *Dao) GetUserByIds(ids []uint32) ([]*UserModel, error) {
	var list []*UserModel
	if err := d.DB.Where("id IN (?) AND re = 0", ids).Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}

// GetUserByEmail get a user by email.
func (d *Dao) GetUserByEmail(email string) (*UserModel, error) {
	u := &UserModel{}
	err := d.DB.Where("email = ? AND re = 0", email).First(u).Error
	if err == gorm.ErrRecordNotFound {
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

	var list []*UserModel

	query := d.DB.Model(&UserModel{}).Where(filter).Offset(int(offset)).Limit(int(limit))

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
	err := d.DB.Where("student_id = ? AND re = 0", studentId).First(u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return u, err
}

func (Dao) RegisterUser(info *RegisterInfo) error {
	user := &UserModel{
		Name:         info.Name,
		Email:        info.Email,
		StudentId:    info.StudentId,
		HashPassword: generatePasswordHash(info.Password),
		Role:         info.Role,
		Re:           false,
	}
	return user.Create()
}
