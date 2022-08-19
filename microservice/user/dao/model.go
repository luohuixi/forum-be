package dao

import (
	"github.com/ShiinaOrez/GoSecurity/security"
)

type UserModel struct {
	Id           uint32 `gorm:"column:id;not null" binding:"required"`
	Name         string `gorm:"column:name;" binding:"required"`
	Email        string `gorm:"column:email;default:null"`
	Avatar       string `gorm:"column:avatar;" binding:"required"`
	StudentId    string `gorm:"column:student_id;"`
	HashPassword string `gorm:"column:hash_password;" binding:"required"`
	Role         string `gorm:"column:role;" binding:"required"`
	Signature    uint32 `gorm:"column:signature;" binding:"required"`
	Re           bool
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
