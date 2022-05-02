package dao

import (
	"forum/model"
	"github.com/jinzhu/gorm"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	DB *gorm.DB
}

// Interface dao
type Interface interface {
	GetUser(id uint32) (*UserModel, error)
	GetUserByIds(ids []uint32) ([]*UserModel, error)
	GetUserByEmail(email string) (*UserModel, error)
	ListUser(offset, limit, lastId uint32, filter *UserModel) ([]*UserModel, error)
	GetUserByStudentId(studentId string) (*UserModel, error)
	RegisterUser(info *RegisterInfo) error
}

// Init init dao
func Init() {
	// lock.Lock()
	// defer lock.Unlock()
	if dao != nil {
		return
	}

	// init db
	model.DB.Init()

	dao = &Dao{
		DB: model.DB.Self,
	}
}

func GetDao() *Dao {
	return dao
}
