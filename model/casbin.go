package model

import (
	"errors"
	"fmt"
	"forum/pkg/constvar"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
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

	a, err := gormadapter.NewAdapter("mysql", config, true) // Your driver and data source.
	if err != nil {
		zap.L().Error("casbin数据库加载失败!", zap.Error(err))
		return nil
	}

	text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub,p.sub) && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		zap.L().Error("字符串加载模型失败!", zap.Error(err))
		return nil
	}
	cb, _ := casbin.NewEnforcer(m, a)

	return cb
}

func (c *Casbin) Init() {
	CB = &Casbin{
		Self: initCasbin(viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.addr"),
			viper.GetString("db.name")),
	}

	rules := [][]string{
		{constvar.Post + ":" + constvar.MuxiPost, constvar.Read},
		{constvar.Post + ":" + constvar.NormalPost, constvar.Read},
	}

	_, err := CB.Self.AddPermissionsForUser(constvar.MuxiRole, rules...)
	if err != nil {
		panic(err)
	}
}

func Enforce(userId uint32, typeName string, data interface{}, act string) (bool, error) {
	object := typeName + ":"

	switch data.(type) {
	case string:
		object += data.(string)
	case uint32:
		object += strconv.Itoa(int(data.(uint32)))
	case int:
		object += strconv.Itoa(data.(int))
	default:
		return false, errors.New("wrong type")
	}

	return CB.Self.Enforce("user:"+strconv.Itoa(int(userId)), object, act)
}

func AddPolicy(userId uint32, typeName string, id uint32, act string) error {
	ok, err := CB.Self.AddPolicy("user:"+strconv.Itoa(int(userId)), typeName+":"+strconv.Itoa(int(id)), act)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("add policy not ok")
	}

	return nil
}

func AddRole(typeName string, id uint32, role string) error {
	ok, err := CB.Self.AddRoleForUser(typeName+":"+strconv.Itoa(int(id)), role)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("add role not ok")
	}

	return nil
}
