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
	// Create() error
	// GetList(uint32) ([]string, error)
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
	defer model.DB.Close()

	dao = &Dao{
		DB: model.DB.Self,
	}
}

func GetDao() *Dao {
	return dao
}
