package dao

import (
	"github.com/ShiinaOrez/GoSecurity/security"
)

type UserModel struct {
	Id           uint32 `json:"id" gorm:"column:id;not null" binding:"required"`
	Name         string `json:"name" gorm:"column:name;" binding:"required"`
	Email        string `json:"email" gorm:"column:email;default:null"`
	Avatar       string `json:"avatar" gorm:"column:avatar;" binding:"required"`
	Role         uint32 `json:"role" gorm:"column:role;" binding:"required"`
	Message      uint32 `json:"message" gorm:"column:message;" binding:"required"`
	HashPassword string `json:"hash_password" gorm:"column:hash_password;" binding:"required"`
	StudentId    string `json:"student_id" gorm:"column:student_id;"`
}

func (u *UserModel) TableName() string {
	return "users"
}

// Create ... create user
func (u *UserModel) Create() error {
	return dao.DB.Create(u).Error
}

// Save ... save user.
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
