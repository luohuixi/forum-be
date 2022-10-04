package dao

import (
	"github.com/ShiinaOrez/GoSecurity/security"
)

type UserModel struct {
	Id           uint32 `gorm:"column:id;not null"`
	Email        string `gorm:"column:email;default:null"`
	Name         string `gorm:"column:name;"`
	Avatar       string `gorm:"column:avatar;"`
	HashPassword string `gorm:"column:hash_password;"`
	Role         string `gorm:"column:role;"`
	Signature    string `gorm:"column:signature;"`
	StudentId    string `gorm:"column:student_id;default:null"`
	Re           bool   `gorm:"column:re;"`
}

func (UserModel) TableName() string {
	return "users"
}

// Create ...
func (u *UserModel) Create() error {
	return dao.DB.Create(u).Error
}

// Save ...
func (u *UserModel) Save() error {
	return dao.DB.Save(u).Error
}

func (u *UserModel) Update() error {
	err := dao.DB.Table("users").Where("id = ?", u.Id).Updates(map[string]interface{}{
		"name":      u.Name,
		"avatar":    u.Avatar,
		"signature": u.Signature,
	}).Error

	return err
}

// generatePasswordHash pwd -> hashPwd
func generatePasswordHash(password string) string {
	return security.GeneratePasswordHash(password)
}

func (u *UserModel) CheckPassword(password string) bool {
	return security.CheckPasswordHash(password, u.HashPassword)
}
