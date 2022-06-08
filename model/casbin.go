package model

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Casbin struct {
	Self *casbin.Enforcer
}

var CB *Casbin

func initCasbin(username, password, addr, DBName string) *casbin.Enforcer {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		username,
		password,
		addr,
		DBName)

	a, _ := gormadapter.NewAdapter("mysql", config, true) // Your driver and data source.

	cb, err := casbin.NewEnforcer("./model.conf", a)
	if err != nil {

	}
	return cb
	// sub := "alice" // the user that wants to access a resource.
	// obj := "data1" // the resource that is going to be accessed.
	// act := "read"  // the operation that the user performs on the resource.
	//
	// if res, _ := e.Enforce(sub, obj, act); res {
	// 	// permit alice to read data1
	// } else {
	// 	// deny the request, show an error
	// }
}

func (c *Casbin) Init() {
	CB = &Casbin{
		Self: initCasbin(viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.addr"),
			viper.GetString("db.name")),
	}
}
