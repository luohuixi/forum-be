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
	GetPostByTime(time string) (*[]PostModel, error)
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
