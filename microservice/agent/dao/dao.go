package dao

import (
	"forum/model"

	"gorm.io/gorm"
)

var dao *Dao

type Dao struct {
	DB *gorm.DB
}

type Interface interface {
	GetValuablePost() (*[]PostModel, error)
	GetPostById(id uint32) (*PostModel, error)
}

func GetDao() *Dao {
	return dao
}

func Init() {
	if dao != nil {
		return
	}

	model.DB.Init()

	dao = &Dao{
		DB: model.DB.Self,
	}
}
