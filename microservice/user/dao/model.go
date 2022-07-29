package dao

import (
	"github.com/ShiinaOrez/GoSecurity/security"
)

type UserModel struct {
	Id           uint32 `json:"id" gorm:"column:id;not null" binding:"required"`
	Name         string `json:"name" gorm:"column:name;" binding:"required"`
	Email        string `json:"email" gorm:"column:email;default:null"`
	Avatar       string `json:"avatar" gorm:"column:avatar;" binding:"required"`
	StudentId    string `json:"student_id" gorm:"column:student_id;"`
	HashPassword string `json:"hash_password" gorm:"column:hash_password;" binding:"required"`
	Role         uint32 `json:"role" gorm:"column:role;" binding:"required"`
	Signature    uint32 `json:"signature" gorm:"column:signature;" binding:"required"`
	Re           bool   `json:"re"`
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

// generatePasswordHash pwd -> hashPwd
func generatePasswordHash(password string) string {
	return security.GeneratePasswordHash(password)
}

func (u *UserModel) CheckPassword(password string) bool {
	return security.CheckPasswordHash(password, u.HashPassword)
}
