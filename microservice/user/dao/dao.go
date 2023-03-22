package dao

import (
	"forum/model"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// Interface dao
type Interface interface {
	GetUser(id uint32) (*UserModel, error)
	GetUserByIds(ids []uint32) ([]*UserModel, error)
	GetUserByEmail(email string) (*UserModel, error)
	ListUser(offset, limit, lastId uint32, filter *UserModel) ([]*UserModel, error)
	GetUserByStudentId(studentId string) (*UserModel, error)
	RegisterUser(info *RegisterInfo) error
	AddPublicPolicy(string, uint32) error

	ListMessage(uint32) ([]string, error)
	CreateMessage(uint32, string) error
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

	// init redis
	model.RedisDB.Init()

	// init casbin
	model.CB.Init()

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
	}
}

func GetDao() *Dao {
	return dao
}
